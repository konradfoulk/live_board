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
				select {
				case client.send <- message:
					// message sent successfully
				default:
					// client not responding, thus is disconnected by default
					// could use a timeout, a buffer, or skip messages to not handle this so harshly
					delete(r.clients, client.username)
				}
			}
		case client := <-r.register:
			r.clients[client.username] = client
			log.Printf("Client %s connected to %s", client.username, r.name)
		case client := <-r.unregister:
			delete(r.clients, client.username)
			log.Printf("Client %s disconnectd from %s", client.username, r.name)
		}
	}
}

func (c *Client) read() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		formattedMsg := c.username + ": " + string(message)
		c.room.broadcast <- []byte(formattedMsg)
	}
}

func (c *Client) write() {
	defer c.conn.Close()

	for message := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
	// when channel closes, send close message
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
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
	fmt.Println("Serving files from ./static")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// upgrade connection to WS
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}

	// get name values
	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		roomName = "general" // default room
	}
	username := r.URL.Query().Get("username")

	// get or create room
	room, exists := hub.rooms[roomName]
	if !exists {
		room = newRoom(roomName, hub)
		room.hub.register <- room
		go room.run()
	}

	client := &Client{
		username: username,
		conn:     conn,
		room:     room,
		send:     make(chan []byte),
	}

	client.room.register <- client
	go client.write()
	go client.read()

}

// var x *type => declares a pointer, x is a memory address which points to another value in memory
// *x => dereferences the pointer, gets the value stored at the memory address x
// &value => gets the memory address for where value is stored
