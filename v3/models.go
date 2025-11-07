package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	username string
	conn     *websocket.Conn
	room     *Room
	hub      *Hub
	send     chan []byte
}

type Room struct {
	name       string
	clients    map[string]*Client
	hub        *Hub
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Hub struct {
	clients          map[string]*Client
	rooms            map[string]*Room
	broadcast        chan []byte
	registerClient   chan *Client
	unregisterClient chan *Client
	registerRoom     chan *Room
	unregisterRoom   chan *Room
}

// read input message from browser and broadcast to room
func (c *Client) read() {
	// remove client from room and hub when client disconnects
	// close WS connection
	defer func() {
		c.room.unregister <- c
		c.hub.unregisterClient <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		msg := WSMessage{
			Type:     "chat",
			Username: c.username,
			Content:  string(message),
		}

		jsonMsg, _ := json.Marshal(msg)
		c.room.broadcast <- jsonMsg
	}
}

// write messages in send channel to browser
func (c *Client) write() {
	// remove client from room and hub when client disconnects
	// close WS connection
	defer func() {
		c.room.unregister <- c
		c.hub.unregisterClient <- c
		c.conn.Close()
	}()

	for message := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
	// send close message when channel closes
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
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
					// client not resopnding, close channel to allow loop to move on
					close(client.send)
				}
			}
		case client := <-r.register:
			r.clients[client.username] = client
			log.Printf("Client %s connected to %s", client.username, r.name)
		case client := <-r.unregister:
			if _, ok := r.clients[client.username]; ok {
				delete(r.clients, client.username)
				log.Printf("Client %s disconnected from %s", client.username, r.name)
			}
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case message := <-h.broadcast:
			for _, client := range h.clients {
				select {
				case client.send <- message:
					// message sent successfully
				default:
					// client not resopnding, close channel to allow loop to move on
					close(client.send)
				}
			}
		case client := <-h.registerClient:
			h.clients[client.username] = client
		case client := <-h.unregisterClient:
			if _, ok := h.clients[client.username]; ok {
				delete(h.clients, client.username)
			}
		case room := <-h.registerRoom:
			h.rooms[room.name] = room
		case room := <-h.unregisterRoom:
			delete(h.rooms, room.name)
		}
	}
}

func makeHub() *Hub {
	return &Hub{
		clients:          make(map[string]*Client),
		rooms:            make(map[string]*Room),
		broadcast:        make(chan []byte),
		registerRoom:     make(chan *Room),
		unregisterRoom:   make(chan *Room),
		registerClient:   make(chan *Client),
		unregisterClient: make(chan *Client),
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

func newClient(username string, room *Room) *Client {
	return &Client{
		username: username,
		room:     room,
		hub:      room.hub,
		send:     make(chan []byte, 256),
	}
}
