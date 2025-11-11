package main

import (
	"encoding/json"
	"log"
	"net/http"

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
	Type     string   `json:"type"`
	Username string   `json:"username,omitempty"`
	Content  string   `json:"content,omitempty"`
	Rooms    []string `json:"rooms,omitempty"`
}

func main() {
	hub := makeHub()
	go hub.run()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API endpoints
	// get rooms (on app load)
	http.HandleFunc("/api/rooms", func(w http.ResponseWriter, r *http.Request) {
		getRooms(hub, w, r)
	})

	// create and delete rooms
	http.HandleFunc("/api/rooms/", func(w http.ResponseWriter, r *http.Request) {
		// handleRooms(hub, w, r)
	})

	// WS endpoint
	// send messages, join and leave instructions, and changes of state (rooms created or deleted, need to update frontend)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// handleWS(hub, w, r)
	})

	// start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getRooms(hub *Hub, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roomNames := []string{}
	for name := range hub.rooms {
		roomNames = append(roomNames, name)
	}

	json.NewEncoder(w).Encode(roomNames)
}
