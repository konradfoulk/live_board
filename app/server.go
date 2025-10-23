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

// type Client struct {
// 	conn *websocket.Conn
// 	hub  *Hub
// 	send chan []byte
// }

// type Hub struct {
// 	clients    map[*Client]bool
// 	broadcast  chan []byte
// 	register   chan *Client
// 	unregister chan *Client
// }

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// add client pointer from channel to hub clients map
			h.clients[client] = true
			log.Println("New client connected")
		case client := <-h.unregister:
			// if the client exists in the hub clients map, delete from the map and close their channel
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("Client disconnected")
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
					// successfully sent
				default:
					// client not responding, thus is disconnected by default
					// could use a timeout, a buffer, or skip messages to not handle this so harshly
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
	// channel was closed, send close message
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}

// func main() {
// 	var hub = newHub()
// 	go hub.run()

// 	// create routes
// 	fs := http.FileServer(http.Dir("./static"))
// 	http.Handle("/", fs)
// 	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// 		handleWebSocket(hub, w, r)
// 	})

// 	fmt.Println("Chat server starting on :8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

func handleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket endpoint hit")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}

	client := &Client{
		conn: conn,
		hub:  hub,
		send: make(chan []byte),
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
