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
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Chat server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head><title>WebSocket Test</title></head>
	<body>
		<h1>Simple WebSocket Test</h1>
		<input type="text" id="messageInput" placeholder="Type a message">
		<button onclick="sendMessage()">Send</button>
		<div id="messages"></div>
		
		<script>
			// Connect to WebSocket
			const ws = new WebSocket('ws://localhost:8080/ws');
			
			ws.onopen = () => {
				console.log('Connected!');
				document.getElementById('messages').innerHTML += '<p>Connected to server!</p>';
			};
			
			ws.onmessage = (event) => {
				document.getElementById('messages').innerHTML += '<p>Server says: ' + event.data + '</p>';
			};
			
			function sendMessage() {
				const input = document.getElementById('messageInput');
				ws.send(input.value);
				document.getElementById('messages').innerHTML += '<p>You: ' + input.value + '</p>';
				input.value = '';
			}
		</script>
	</body>
	</html>`

	fmt.Fprintf(w, html)
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
