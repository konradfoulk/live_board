package main

import "github.com/gorilla/websocket"

type Client struct {
	// username string
	conn *websocket.Conn
	room *Room
	send chan []byte
}

type Room struct {
	name string
	// clients    map[string]*Client
	clients    map[*Client]bool
	hub        *Hub
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Hub struct {
	rooms      map[string]*Room
	register   chan *Room
	unregister chan *Room
}
