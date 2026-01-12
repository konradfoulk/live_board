# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Real-time multi-room chat application built with Go backend and vanilla JavaScript frontend, using WebSockets for live communication.

## Build & Run Commands

```bash
# Build the server
go build

# Run the server (serves on port 8080)
go run server.go models.go db.go
```

The server serves static files from `./static/` and creates/initializes the SQLite database (`app.db`) on startup using `schema.sql`.

## Architecture

### Backend (Go)

**Hub Pattern** - Central broker managing all clients and rooms:
- `Hub` (models.go) - Maintains maps of all connected clients and active rooms, runs its own goroutine for coordination
- `Client` (models.go) - Represents a WebSocket connection with dedicated read/write goroutines
- `Room` (models.go) - Chat room with broadcast channel, distributes messages to members

**Concurrency Model:**
- Each Client spawns read and write goroutines
- Each Room runs a goroutine listening on its broadcast channel
- RWMutex protects shared maps (clients, rooms)
- Channels used for all inter-goroutine communication

**Server Endpoints (server.go):**
- `GET /` - Static file server
- `POST /api/rooms` - Create room
- `POST /api/rooms/{roomName}` - Delete room
- `WS /ws?username={username}&password={password}` - WebSocket connection with bcrypt authentication

### Frontend (static/)

- `index.html` - Main page with join modal and chat interface
- `scripts.js` - WebSocket client, message handling, chat logic
- `ui.js` - DOM manipulation, room buttons, form events

### Database (SQLite)

Schema in `schema.sql`: users, rooms, messages tables. Messages foreign-key to rooms (with cascade delete) and users.

## WebSocket Message Types

- `join_room` - Client joins a room
- `message` - Chat message (subtypes: chat_message, join_message, leave_message, init_chat)
- `create_room` / `delete_room` - Room lifecycle broadcasts
- `init_rooms` - Initial room list for new connections
- `user_count` - Global connected user count

## Key Implementation Details

- Users auto-created on first login (bcrypt password hashing)
- Last 200 messages loaded when joining a room
- Message persistence is asynchronous via goroutine
- 256-message buffer per client channel
- "general" room created by default on startup
