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

type WSMessage struct {
	Type  string   `json:"type"`
	Room  string   `json:"room,omitempty"`
	Rooms []string `json:"rooms,omitempty"`
}

func main() {
	// make hub
	hub := makeHub()
	go hub.run()

	// serve app
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API endpoints
	// create room
	http.HandleFunc("/api/rooms", func(w http.ResponseWriter, r *http.Request) {
		createRoom(hub, w, r)
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

	// make and start room
	room := newRoom(roomName)
	hub.registerRoom <- room

	roomCreated := <-hub.createRoom // make sure client doesn't get room that doesn't exist yet

	// push update to frontend clients
	msg := WSMessage{
		Type: "create_room",
		Room: roomCreated,
	}
	jsonMsg, _ := json.Marshal(msg)
	hub.broadcast <- jsonMsg

	// send success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"name": roomName})
}

func handleWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket endpoint hit")

	// upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade failed:", err)
		return
	}

	username := r.URL.Query().Get("username")
	client := newClient(username, conn)

	hub.registerClient <- client
	go client.write()

	// send initial state to frontend
	rooms := <-hub.initRooms
	msg := WSMessage{
		Type:  "init_rooms",
		Rooms: rooms,
	}
	jsonMsg, _ := json.Marshal(msg)
	client.send <- jsonMsg
}
