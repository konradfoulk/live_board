package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WS upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// global db variable for access
var db *sql.DB

func main() {
	// initialize database
	db = initDatabase()

	// make hub
	hub := makeHub()
	go hub.run()

	// make default room
	db.Exec("INSERT OR IGNORE INTO rooms (name) VALUES ('general')")
	defaultRoom := newRoom("general")
	go defaultRoom.run()
	hub.rooms[defaultRoom.name] = defaultRoom

	// serve app
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API endpoints
	// create room
	http.HandleFunc("/api/rooms", func(w http.ResponseWriter, r *http.Request) {
		createRoom(hub, w, r)
	})

	// delete room
	http.HandleFunc("/api/rooms/", func(w http.ResponseWriter, r *http.Request) {
		deleteRoom(hub, w, r)
	})

	// WS endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWS(hub, w, r)
	})

	// start server
	log.Println("chat server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createRoom(hub *Hub, w http.ResponseWriter, r *http.Request) {

}
