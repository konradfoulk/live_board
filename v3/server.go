package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// WS upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WS message format
type WSMessage struct {
	Type     string   `json:"type"`
	Username string   `json:"username,omitempty"`
	Content  string   `json:"content,omitempty"`
	Rooms    []string `json:"rooms,omitempty"`
}
