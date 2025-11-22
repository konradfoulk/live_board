package main

import (
	"log"
	"slices"
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
	name       string
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
}

type Hub struct {
	clients          map[string]*Client
	clientsMutex     sync.RWMutex
	rooms            map[string]*Room
	roomsList        []string // for order and state
	roomsMutex       sync.RWMutex
	broadcast        chan []byte
	registerClient   chan *Client
	unregisterClient chan *Client
	registerRoom     chan *Room
	unregisterRoom   chan *Room
	initRooms        chan []string
	createRoom       chan string
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
				c.room.unregister <- c
			}

			room := c.hub.rooms[msg.Room]
			c.room = room
			room.register <- c
		}
	}

	// receive room join request from front end
	// unregister from current room (if not room === "")
	// register for the new room
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			log.Printf("%s joined %s", client.username, r.name)
		case client := <-r.unregister:
			log.Printf("%s left %s", client.username, r.name)
		}
	}
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
		case client := <-h.registerClient:
			h.roomsMutex.RLock()

			h.clientsMutex.Lock()
			h.clients[client.username] = client
			h.clientsMutex.Unlock()

			h.initRooms <- h.roomsList
			h.roomsMutex.RUnlock()

			log.Printf("%s connected to hub", client.username)
		case room := <-h.registerRoom:
			h.roomsMutex.Lock()
			h.rooms[room.name] = room
			h.roomsList = append(h.roomsList, room.name)

			h.createRoom <- room.name
			h.roomsMutex.Unlock()

			log.Printf("created room %s", room.name)
		case room := <-h.unregisterRoom:
			h.roomsMutex.Lock()
			delete(h.rooms, room.name)
			h.roomsList = slices.DeleteFunc(h.roomsList, func(name string) bool {
				return name == room.name
			})
			h.roomsMutex.Unlock()

			log.Printf("deleted room %s", room.name)
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
		name:       name,
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func makeHub() *Hub {
	return &Hub{
		clients:          make(map[string]*Client),
		rooms:            make(map[string]*Room),
		roomsList:        []string{},
		broadcast:        make(chan []byte),
		registerRoom:     make(chan *Room),
		unregisterRoom:   make(chan *Room),
		registerClient:   make(chan *Client),
		unregisterClient: make(chan *Client),
		initRooms:        make(chan []string),
		createRoom:       make(chan string),
	}
}
