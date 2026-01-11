package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

// global db variable for access
var db *sql.DB

// WS upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	Type        string   `json:"type"`
	MessageType string   `json:"messageType,omitempty"`
	Username    string   `json:"username,omitempty"`
	Room        string   `json:"room,omitempty"`
	Content     string   `json:"content,omitempty"`
	Rooms       []string `json:"rooms,omitempty"`
}

func main() {
	// initialize database
	db = initDatabase()

	// make hub
	hub := makeHub()
	go hub.run()

	// make default room
	db.Exec("INSERT OR IGNORE INTO rooms (name) VALUES ('general')")

	// load state from db in memory
	rows, _ := db.Query("SELECT name FROM rooms")
	defer rows.Close()
	for rows.Next() {
		var roomName string
		rows.Scan(&roomName)

		room := newRoom(roomName)
		go room.run()
		hub.rooms[room.name] = room
	}

	// serve app
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API endpoints
	// create room
	http.HandleFunc("/api/rooms", func(w http.ResponseWriter, r *http.Request) {
		createRoom(hub, w, r)
	})

	// delete room
	http.HandleFunc("/api/rooms/", func(w http.ResponseWriter, r *http.Request) {
		deleteRoom(hub, w, r)
	})

	// WS endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWS(hub, w, r)
	})

	// start server
	log.Println("chat server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createRoom(hub *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// cast request into Go struct and get room name from request body
	var reqBody struct {
		Name string `json:"name"`
	}
	json.NewDecoder(r.Body).Decode(&reqBody)

	roomName := reqBody.Name

	// persist, make, and start room and register with hub
	hub.roomsMutex.Lock()
	if _, err := db.Exec("INSERT INTO rooms (name) VALUES (?)", roomName); err != nil {
		hub.roomsMutex.Unlock()

		// Send error response to client
		w.WriteHeader(http.StatusConflict) // 409
		json.NewEncoder(w).Encode(map[string]string{
			"error": "room already exists",
		})
		return
	}

	room := newRoom(roomName)
	go room.run()

	hub.rooms[room.name] = room

	log.Printf("created room %s", room.name)
	hub.roomsMutex.Unlock()

	// push update to frontend clients
	msg := WSMessage{
		Type: "create_room",
		Room: room.name,
	}
	jsonMsg, _ := json.Marshal(msg)
	hub.broadcast <- jsonMsg

	// send success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"name": roomName})
}

func deleteRoom(hub *Hub, w http.ResponseWriter, r *http.Request) {
	roomName := strings.TrimPrefix(r.URL.Path, "/api/rooms/")

	// delete on backend
	hub.roomsMutex.Lock()
	room := hub.rooms[roomName]
	if room != nil {
		// delete from db
		if _, err := db.Exec("DELETE FROM rooms WHERE name = ?", room.name); err != nil {
			hub.roomsMutex.Unlock()

			// send error response
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to delete room",
			})
			return
		}

		// push update to frontend clients
		msg := WSMessage{
			Type: "delete_room",
			Room: room.name,
		}
		jsonMsg, _ := json.Marshal(msg)
		hub.broadcast <- jsonMsg

		// remove clients from room
		room.clientsMutex.Lock()
		for _, client := range room.clients {
			client.room = nil
			delete(room.clients, client.username)

			log.Printf("%s left %s", client.username, room.name)
		}
		room.clientsMutex.Unlock()

		// delete from memory
		delete(hub.rooms, room.name)

		close(room.broadcast)

		log.Printf("deleted room %s", room.name)
	}
	hub.roomsMutex.Unlock()

	// send success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"name": roomName})
}

// need better error handling here for the client side
func handleWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket endpoint hit")

	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	// check if user exists
	var userID int
	var storedHash string
	err := db.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", username).Scan(&userID, &storedHash)

	// authenticate
	if err == sql.ErrNoRows {
		// new user
		// hash password and create field
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		result, _ := db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, passwordHash)
		id, _ := result.LastInsertId()
		userID = int(id)
		log.Printf("created new user: %s (id: %d)", username, userID)
	} else if err != nil {
		// database error
		log.Println("database error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		// authenticate existing user
		if err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	log.Printf("%s (id: %d) connecting", username, userID)

	// upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade failed:", err)
		return
	}

	client := newClient(username, conn, hub)
	go client.write()
	go client.read()

	hub.roomsMutex.RLock()
	hub.clientsMutex.Lock()
	hub.clients[client.username] = client

	log.Printf("%s connected to hub", client.username)
	hub.clientsMutex.Unlock()

	// get list of existing rooms
	rows, _ := db.Query("SELECT name FROM rooms ORDER BY created_at ASC")
	defer rows.Close()
	roomsList := []string{}
	for rows.Next() {
		var roomName string
		rows.Scan(&roomName)
		roomsList = append(roomsList, roomName)
	}

	// send initial state to frontend
	msg := WSMessage{
		Type:  "init_rooms",
		Rooms: roomsList,
	}
	jsonMsg, _ := json.Marshal(msg)
	client.send <- jsonMsg
	hub.roomsMutex.RUnlock()
}
