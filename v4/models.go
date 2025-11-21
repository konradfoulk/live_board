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
}

type Room struct {
	name string
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
		}
	}
}

func newClient(username string, conn *websocket.Conn) *Client {
	return &Client{
		username: username,
		conn:     conn,
		send:     make(chan []byte, 256),
	}
}

func newRoom(name string) *Room {
	return &Room{
		name: name,
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
