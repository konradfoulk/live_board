// Lesson 004: Understanding WebSockets in Go
// ============================================
// WebSockets enable real-time, bidirectional communication between browser and server

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader - this upgrades HTTP connections to WebSocket connections
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin returns true to allow connections from any origin
	// In production, you'd want to be more restrictive
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents a connected user
type Client struct {
	conn *websocket.Conn // The WebSocket connection
	send chan []byte     // Channel to send messages to this client
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool // Map of connected clients
	broadcast  chan []byte      // Channel for messages to broadcast
	register   chan *Client     // Channel to register new clients
	unregister chan *Client     // Channel to unregister clients
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// run handles all hub events using channels and select
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// New client connected
			h.clients[client] = true
			log.Println("Client connected. Total clients:", len(h.clients))

		case client := <-h.unregister:
			// Client disconnected
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("Client disconnected. Total clients:", len(h.clients))
			}

		case message := <-h.broadcast:
			// Broadcast message to all clients
			for client := range h.clients {
				select {
				case client.send <- message:
					// Successfully sent
				default:
					// Client's send channel is full, close it
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// HOW WEBSOCKETS WORK:
// 1. Browser makes HTTP request to /ws endpoint
// 2. Server "upgrades" the connection from HTTP to WebSocket
// 3. Now both sides can send messages at any time (not request/response)
// 4. Messages flow through channels to handle concurrency safely

// KEY CONCEPTS:
// - Upgrader: Converts HTTP connection to WebSocket
// - Client: Represents each connected user with their own goroutine
// - Hub: Central manager that broadcasts messages to all clients
// - Channels: Safe communication between goroutines

// FLOW:
// User connects -> Create Client -> Register with Hub -> Listen for messages
// User sends message -> Hub receives -> Hub broadcasts to all clients
// User disconnects -> Unregister from Hub -> Clean up resources

func main() {
	fmt.Println("This is a lesson file showing WebSocket concepts")
	fmt.Println("The actual implementation will be in your server.go")
	fmt.Println("\nKey components:")
	fmt.Println("1. Upgrader - converts HTTP to WebSocket")
	fmt.Println("2. Client - represents a connected user")
	fmt.Println("3. Hub - manages all clients and broadcasts")
	fmt.Println("4. Channels - safe concurrent communication")
}

// Next: We'll implement these concepts in your server.go!
