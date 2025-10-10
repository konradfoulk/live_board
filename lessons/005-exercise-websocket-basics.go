// Exercise 005: WebSocket Step by Step
// ======================================
// Let's build this piece by piece so it makes sense!

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// STEP 1: The Upgrader
// This is like a doorman that converts regular HTTP visitors into WebSocket VIPs
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (for development)
	},
}

// STEP 2: Simple WebSocket Handler (no Hub yet!)
func simpleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 1. Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	fmt.Println("🎉 Client connected!")

	// 2. This is where the GOROUTINE would normally be!
	// For now, we'll just echo messages back
	for {
		// Read message from browser
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected")
			break
		}

		fmt.Printf("Received: %s\n", message)

		// Echo the message back
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println("Write failed:", err)
			break
		}
	}
}

// WHERE ARE THE GOROUTINES?
// Great question! They come in when we have MULTIPLE clients:
//
// Without goroutines (BAD):
// Client 1 connects -> Server handles Client 1
// Client 2 tries to connect -> WAITS until Client 1 is done!
//
// With goroutines (GOOD):
// Client 1 connects -> go handleClient1()
// Client 2 connects -> go handleClient2()
// Both run simultaneously!

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", simpleWebSocket)

	fmt.Println("Server starting on :8080")
	fmt.Println("Try connecting with a WebSocket client!")
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

// TODO for you:
// 1. Run this server
// 2. Open http://localhost:8080 in your browser
// 3. Type messages and see them echo back
// 4. Open ANOTHER browser tab to the same URL
// 5. Notice the problem: Can you chat between tabs? (Spoiler: No!)
//
// This is why we need the Hub pattern - to connect multiple clients!
