package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	username       string
	conn           *websocket.Conn
	send           chan []byte
	disconnectOnce sync.Once
	room           *Room
	hub            *Hub
}

type Room struct {
	name         string
	clients      map[string]*Client
	clientsMutex sync.RWMutex
	broadcast    chan []byte
}

type Hub struct {
	clients      map[string]*Client
	clientsMutex sync.RWMutex
	rooms        map[string]*Room
	roomsMutex   sync.RWMutex
	broadcast    chan []byte
}

func (c *Client) write() {
	defer c.disconnectOnce.Do(c.disconnect)

	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	// channel closed, send close message
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func (c *Client) read() {
	defer c.disconnectOnce.Do(c.disconnect)

	for {
		var msg WSMessage
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			// connection was closed, client left
			break
		}

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

				log.Printf("%s joined %s", c.username, room.name)

				// send presence update
				formattedMsg := WSMessage{
					Type:        "message",
					MessageType: "join_message",
					Username:    c.username,
					Room:        room.name,
				}
				jsonMsg, _ := json.Marshal(formattedMsg)
				room.broadcast <- jsonMsg

				room.clientsMutex.Unlock()
			}
			c.hub.roomsMutex.RUnlock()
		case "message":
			if c.room != nil && msg.Room == c.room.name {
				// format and forward message to room broadcast
				formattedMsg := WSMessage{
					Type:        "message",
					MessageType: "chat_message",
					Username:    c.username,
					Room:        c.room.name,
					Content:     msg.Content,
				}
				jsonMsg, _ := json.Marshal(formattedMsg)
				c.room.broadcast <- jsonMsg
			}
		}
	}
}

func (c *Client) disconnect() {
	if c.room != nil {
		c.room.unregister(c)
	}
	c.hub.unregister(c)

	c.conn.Close()

	close(c.send)
}

func (r *Room) unregister(client *Client) {
	r.clientsMutex.Lock()
	if client.room != nil {
		client.room = nil
		delete(r.clients, client.username)

		log.Printf("%s left %s", client.username, r.name)

		// send presence update
		msg := WSMessage{
			Type:        "message",
			MessageType: "leave_message",
			Username:    client.username,
			Room:        r.name,
		}
		jsonMsg, _ := json.Marshal(msg)
		r.broadcast <- jsonMsg
	}
	r.clientsMutex.Unlock()
}

func (r *Room) run() {
	for message := range r.broadcast {
		r.clientsMutex.RLock()
		for _, client := range r.clients {
			select {
			case client.send <- message:
				// message sent successfully
			default:
				log.Printf("%s not responding", client.username)
				client.conn.Close()
			}
		}
		r.clientsMutex.RUnlock()
	}
}

func (h *Hub) unregister(client *Client) {
	h.clientsMutex.Lock()
	delete(h.clients, client.username)

	log.Printf("%s disconnected from hub", client.username)
	h.clientsMutex.Unlock()
}

func (h *Hub) run() {
	for message := range h.broadcast {
		h.clientsMutex.RLock()
		for _, client := range h.clients {
			select {
			case client.send <- message:
				// message sent successfully
			default:
				log.Printf("%s not responding", client.username)
				client.conn.Close()
			}
		}
		h.clientsMutex.RUnlock()
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
		name:      name,
		clients:   make(map[string]*Client),
		broadcast: make(chan []byte),
	}
}

func makeHub() *Hub {
	return &Hub{
		clients:   make(map[string]*Client),
		rooms:     make(map[string]*Room),
		broadcast: make(chan []byte),
	}
}
