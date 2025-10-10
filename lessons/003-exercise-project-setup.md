# Exercise 003: Setting Up Your WebSocket Chat Project

## Your Task:
Set up the project dependencies and create the basic server structure.

## Steps to Complete:

### 1. First, update your go.mod file to include the Gorilla WebSocket package
Run this command in your terminal (in the app directory):
```bash
go get github.com/gorilla/websocket
```

### 2. Create a basic HTTP server file
Create a new file called `server.go` with this starter code that you need to complete:

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    // TODO: Add the websocket import here
)

func main() {
    // TODO: Set up a route for the home page "/"
    // Hint: Use http.HandleFunc()
    
    // TODO: Set up a route for websocket connections "/ws"
    
    fmt.Println("Chat server starting on :8080")
    // TODO: Start the server on port 8080
    // Hint: Use log.Fatal(http.ListenAndServe(...))
}

func homePage(w http.ResponseWriter, r *http.Request) {
    // TODO: Write a simple HTML response
    // For now, just send "Chat Server Home Page"
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    // We'll implement this next
    fmt.Println("WebSocket endpoint hit")
}
```

### 3. Test Your Setup
Once you've completed the TODOs:
1. Run `go run server.go`
2. Open your browser to http://localhost:8080
3. You should see your home page message

## Need a Hint?
- For imports: You need `github.com/gorilla/websocket`
- For routes: `http.HandleFunc("/", homePage)`
- For server: `log.Fatal(http.ListenAndServe(":8080", nil))`

When you're done, respond with "Done" or "I need a Hint" if you get stuck!
