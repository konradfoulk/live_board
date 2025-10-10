// Lesson 002: Goroutines and Channels for Real-time Chat
// ========================================================
// Konrad, these are Go's superpowers for handling multiple chat users simultaneously!

package main

import (
	"fmt"
	"time"
)

// GOROUTINES - Lightweight threads for concurrent execution
// Think of them like async functions in JavaScript, but much more efficient

// Simulating a user sending messages
func userSendMessages(userName string, messageChannel chan string) {
	messages := []string{
		fmt.Sprintf("%s: Hey everyone!", userName),
		fmt.Sprintf("%s: Building something cool today", userName),
		fmt.Sprintf("%s: Go is pretty awesome!", userName),
	}

	for _, msg := range messages {
		time.Sleep(1 * time.Second) // Simulate typing delay
		messageChannel <- msg       // Send message to channel
		fmt.Printf("[SENT] %s\n", msg)
	}
}

// CHANNELS - Go's way of communication between goroutines
// Think of them as pipes that goroutines use to pass data safely

func chatServer() {
	// Create a channel for string messages
	// This is like a message queue that goroutines can send to and receive from
	messageChannel := make(chan string)

	// Start multiple users sending messages concurrently
	// The 'go' keyword launches each function in its own goroutine
	go userSendMessages("Konrad", messageChannel)
	go userSendMessages("Alex", messageChannel)

	// Receive and display messages from all users
	// This will run 6 times (3 messages × 2 users)
	for i := 0; i < 6; i++ {
		receivedMsg := <-messageChannel // Receive from channel
		fmt.Printf("[BROADCAST] %s\n", receivedMsg)
	}
}

// BUFFERED CHANNELS - Channels with a capacity
// Like having a message queue with a maximum size

func bufferedExample() {
	// Create a buffered channel with capacity of 3
	// Can hold up to 3 messages before blocking
	notifications := make(chan string, 3)

	// These won't block because buffer has space
	notifications <- "User joined the chat"
	notifications <- "New message received"
	notifications <- "User is typing..."

	// Read all notifications
	fmt.Println(<-notifications)
	fmt.Println(<-notifications)
	fmt.Println(<-notifications)
}

// SELECT STATEMENT - Handle multiple channels
// Like a switch statement for channels

func selectExample() {
	messages := make(chan string)
	signals := make(chan bool)

	// Simulate incoming messages
	go func() {
		time.Sleep(2 * time.Second)
		messages <- "Hello from goroutine!"
	}()

	// Simulate a timeout signal
	go func() {
		time.Sleep(3 * time.Second)
		signals <- true
	}()

	// Use select to handle whichever channel receives data first
	for i := 0; i < 2; i++ {
		select {
		case msg := <-messages:
			fmt.Println("Received message:", msg)
		case sig := <-signals:
			fmt.Println("Received signal:", sig)
		case <-time.After(1 * time.Second):
			fmt.Println("Timeout: No activity for 1 second")
		}
	}
}

func main() {
	fmt.Println("=== Goroutines & Channels Lesson ===")

	fmt.Println("1. Basic Chat Server Simulation:")
	fmt.Println("---------------------------------")
	chatServer()

	fmt.Println("\n2. Buffered Channel Example:")
	fmt.Println("-----------------------------")
	bufferedExample()

	fmt.Println("\n3. Select Statement Example:")
	fmt.Println("-----------------------------")
	selectExample()
}

// Key Concepts for Our Chat App:
// - Each websocket connection will run in its own goroutine
// - Channels will pass messages between connections
// - Select statements will handle multiple events (messages, disconnects, etc.)
//
// Run this to see concurrent message handling in action!
