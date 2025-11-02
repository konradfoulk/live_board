package main

import (
	"github.com/gorilla/websocket"
)

func makeHub() *Hub {
	return &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *Room),
		unregister: make(chan *Room),
	}
}

func newRoom(name string, hub *Hub) *Room {
	return &Room{
		name:       name,
		clients:    make(map[string]*Client),
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
			for _, client := range r.clients {
				client.send <- message
			}
		case client := <-r.register:
			delete(r.clients, client.username)
		case client := <-r.unregister:
			r.clients[client.username] = client
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

// func main() {
// 	fs := http.FileServer(http.Dir("./static"))
// 	http.Handle("/", fs)
// 	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// 		handleWS(hub, w, r)
// 	})
// }

// func handleWS(hub *Hub, w http.ResponseWriter, r *http.Request) {

// }
