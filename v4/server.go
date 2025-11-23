package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/gorilla/websocket"
)

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
	Rooms       []string `json:"rooms,omitempty"`
}

func main() {
	// make hub
	hub := makeHub()
	go hub.run()
	// make default room
	defaultRoom := newRoom("general")
	go defaultRoom.run()
	hub.rooms[defaultRoom.name] = defaultRoom
	hub.roomsList = append(hub.roomsList, defaultRoom.name)

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

	// make and start room
	room := newRoom(roomName)
	go room.run()

	// register room with hub
	hub.roomsMutex.Lock()
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

func deleteRoom(hub *Hub, w http.ResponseWriter, r *http.Request) {
	roomName := strings.TrimPrefix(r.URL.Path, "/api/rooms/")

	hub.roomsMutex.RLock()
	room := hub.rooms[roomName]
	hub.roomsMutex.RUnlock()

	// push update to frontend clients
	msg := WSMessage{
		Type: "delete_room",
		Room: room.name,
	}
	jsonMsg, _ := json.Marshal(msg)
	hub.broadcast <- jsonMsg

	// delete on backend
	hub.roomsMutex.Lock()
	if hub.rooms[room.name] != nil {
		room.clientsMutex.Lock()
		for _, client := range room.clients {
			client.room = nil
			delete(room.clients, client.username)

			// room.broadcast <- client left this room

			log.Printf("%s left %s", client.username, room.name)
		}
		room.clientsMutex.Unlock()

		delete(hub.rooms, room.name)
		hub.roomsList = slices.DeleteFunc(hub.roomsList, func(name string) bool {
			return name == room.name
		})

		log.Printf("deleted room %s", room.name)
	}
	hub.roomsMutex.Unlock()

	// send success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"name": roomName}) // how to get roomName from url (api endpoint)
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
	client := newClient(username, conn, hub)
	go client.write()
	go client.read()

	hub.roomsMutex.RLock()
	hub.clientsMutex.Lock()
	hub.clients[client.username] = client

	log.Printf("%s connected to hub", client.username)
	hub.clientsMutex.Unlock()

	// send initial state to frontend
	msg := WSMessage{
		Type:  "init_rooms",
		Rooms: hub.roomsList,
	}
	hub.roomsMutex.RUnlock()
	jsonMsg, _ := json.Marshal(msg)
	client.send <- jsonMsg
}
