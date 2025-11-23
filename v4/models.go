package main

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	username string
	conn     *websocket.Conn
	send     chan []byte
	room     *Room
	hub      *Hub
}

type Room struct {
	name         string
	clients      map[string]*Client
	clientsMutex sync.RWMutex
	broadcast    chan []byte
	// unregister   chan *Client
}

type Hub struct {
	clients          map[string]*Client
	clientsMutex     sync.RWMutex
	rooms            map[string]*Room
	roomsList        []string // for order and state
	roomsMutex       sync.RWMutex
	broadcast        chan []byte
	unregisterClient chan *Client
}

func (c *Client) write() {
	for message := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (c *Client) read() {
	for {
		var msg WSMessage
		c.conn.ReadJSON(&msg) // break if there is an error (means connection was closed, ie client left)

		switch msg.Type {
		case "join_room":
			if c.room != nil {
				c.room.unregister(c)
			}

			c.hub.roomsMutex.RLock()
			if room := c.hub.rooms[msg.Room]; room != nil {
				room.clientsMutex.Lock()
				c.room = room
				room.clients[c.username] = c

				// room.broadcast <- client joined this room

				log.Printf("%s joined %s", c.username, room.name)
				room.clientsMutex.Unlock()
			}
			c.hub.roomsMutex.RUnlock()
		case "message":
			if c.room != nil {

			}
		}
	}

	// receive room join request from front end
	// unregister from current room (if not room === "")
	// register for the new room
}

func (r *Room) unregister(client *Client) {
	r.clientsMutex.Lock()
	client.room = nil
	delete(r.clients, client.username)

	// room.broadcast <- client left this room

	log.Printf("%s left %s", client.username, r.name)
	r.clientsMutex.Unlock()
}

func (r *Room) run() {
	// for {
	// 	select {
	// 	case client := <-r.unregister:
	// 		client.room = nil

	// 		delete(r.clients, client.username)

	// 		log.Printf("%s left %s", client.username, r.name)
	// 	}
	// }
}

func (h *Hub) run() {
	for {
		select {
		case message := <-h.broadcast:
			h.clientsMutex.RLock()
			for _, client := range h.clients {
				select {
				case client.send <- message:
					// message sent successfully
				default:
					log.Printf("%s not responding", client.username)
					close(client.send)
				}
			}
			h.clientsMutex.RUnlock()
		}
	}
}

func newClient(username string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		username: username,
		conn:     conn,
		send:     make(chan []byte, 256),
		hub:      hub,
	}
}

func newRoom(name string) *Room {
	return &Room{
		name:    name,
		clients: make(map[string]*Client),
		// unregister: make(chan *Client),
	}
}

func makeHub() *Hub {
	return &Hub{
		clients:          make(map[string]*Client),
		rooms:            make(map[string]*Room),
		roomsList:        []string{},
		broadcast:        make(chan []byte),
		unregisterClient: make(chan *Client),
	}
}
