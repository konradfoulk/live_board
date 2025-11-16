package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// WS upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WS message format
type WSMessage struct {
	Type     string      `json:"type"` // chat, get rooms, get chats, join/leave room, new message (notifaction, impliment as a counter or ignore if user is in that room on front end)
	Username string      `json:"username,omitempty"`
	Content  string      `json:"content,omitempty"`
	Rooms    []string    `json:"rooms,omitempty"`
	Messages []WSMessage `json:"messages,omitempty"`
	Room     string      `json:"room,omitempty"`
}

func main() {
	// make hub
	hub := makeHub()
	go hub.run()
	// create default room
	general := newRoom("general", hub)
	hub.registerRoom <- general
	go general.run()

	// serve app
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API endpoints
	// create a new room
	http.HandleFunc("/api/rooms", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getRooms(hub, w, r) // get rid of get rooms, moving operation to websockets
		case "POST":
			createRoom(hub, w, r)
		}
	})

	// delete rooms
	http.HandleFunc("/api/rooms/", func(w http.ResponseWriter, r *http.Request) {
		deleteRoom(hub, w, r)
	})

	// WS endpoint
	// send messages, join and leave instructions, and changes of state (rooms created or deleted, need to update frontend)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWS(hub, w, r)
	})

	// start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// get rid of this, send intial state as first request via websocket
func getRooms(hub *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roomNames := []string{}
	hub.roomsMutex.RLock()
	for name := range hub.rooms {
		roomNames = append(roomNames, name)
	}
	hub.roomsMutex.RUnlock()

	json.NewEncoder(w).Encode(roomNames)
}

func createRoom(hub *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get room name from request body
	var reqBody struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	roomName := reqBody.Name

	// make and start room
	room := newRoom(roomName, hub)
	hub.registerRoom <- room
	go room.run()

	// broadcast state change to clients

	// send success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"name": roomName})
}

func deleteRoom(hub *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roomName := strings.TrimPrefix(r.URL.Path, "/api/rooms/")

	// might be unnecessary
	hub.roomsMutex.RLock()
	room, exists := hub.rooms[roomName]
	hub.roomsMutex.RUnlock()
	if !exists {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	// send clients in room to general (unregister from room)
	// broadcast state change to clients? (just need to send them updated list of roomnames so they can create or delete buttons with appropriate functions)
	// close room channels? (stop room go routine)

	hub.unregisterRoom <- room

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "deleted " + roomName})
}

func handleWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade failed:", err)
		return
	}

	// send initial room state through ws

	// create new client and add them to app
	// add them to the first room in the app, if there is no room, perhaps don't add them to a room?
	username := r.URL.Query().Get("username")
	hub.roomsMutex.RLock()
	room := hub.rooms["general"]
	hub.roomsMutex.RUnlock()

	client := newClient(username, conn, room)

	room.register <- client
	hub.registerClient <- client
	go client.write()
	go client.read()
}
