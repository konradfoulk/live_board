package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ws upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// create routes
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Chat server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket endpoint hit")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	fmt.Println("🎉 Client connected!")

	for {
		// read message from browser
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected")
			break
		}

		fmt.Printf("Received: %s\n", message)

		// echo message back
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println("Write failed:", err)
			break
		}
	}
}
