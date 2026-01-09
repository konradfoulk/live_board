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

func main() {

}
