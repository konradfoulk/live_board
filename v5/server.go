package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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
	hub.roomsList = append(hub.roomsList, room.name)

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
