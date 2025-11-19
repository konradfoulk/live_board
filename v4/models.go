package main

import (
	"encoding/json"
	"log"

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
	rooms            map[string]*Room
	broadcast        chan []byte
	registerClient   chan *Client
	unregisterClient chan *Client
	registerRoom     chan *Room
	unregisterRoom   chan *Room
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

			log.Printf("%s connected to hub", client.username)
		case room := <-h.registerRoom:
			h.rooms[room.name] = room

			// push update to frontend
			msg := WSMessage{
				Type: "create_room",
				Room: room.name,
			}
			jsonMsg, _ := json.Marshal(msg)
			h.broadcast <- jsonMsg
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
		broadcast:        make(chan []byte),
		registerRoom:     make(chan *Room),
		unregisterRoom:   make(chan *Room),
		registerClient:   make(chan *Client),
		unregisterClient: make(chan *Client),
	}
}
