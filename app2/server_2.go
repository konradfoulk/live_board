package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ws upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func makeHub() *Hub {
	return &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *Room),
		unregister: make(chan *Room),
	}
}

func newRoom(name string, hub *Hub) *Room {
	return &Room{
		name: name,
		// clients:    make(map[string]*Client),
		clients:    make(map[*Client]bool),
		hub:        hub,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case room := <-h.register:
			h.rooms[room.name] = room
		case room := <-h.unregister:
			delete(h.rooms, room.name)
		}
	}
}

func (r *Room) run() {
	for {
		select {
		case message := <-r.broadcast:
			// for _, client := range r.clients {
			// 	client.send <- message
			// }
			for client := range r.clients {
				client.send <- message
			}
		case client := <-r.register:
			// r.clients[client.username] = client
			r.clients[client] = true
		case client := <-r.unregister:
			// delete(r.clients, client.username)
			delete(r.clients, client)
		}
	}
}

func (c *Client) read() {
	for {
		_, message, _ := c.conn.ReadMessage()

		c.room.broadcast <- message
	}
}

func (c *Client) write() {
	for message := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
}

func main() {
	var hub = makeHub()
	go hub.run()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWS(hub, w, r)
	})

	fmt.Println("Chat server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// upgrade connection to WS
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}

	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		roomName = "general" // default room
	}

	// get or create room
	room, exists := hub.rooms[roomName]
	if !exists {
		room = newRoom(roomName, hub)
		room.hub.register <- room
		go room.run()
	}

	client := &Client{
		conn: conn,
		room: room,
		send: make(chan []byte),
	}

	client.room.register <- client
	go client.write()
	go client.read()

}

// var x *type => declares a pointer, x is a memory address which points to another value in memory
// *x => dereferences the pointer, gets the value stored at the memory address x
// &value => gets the memory address for where value is stored
