# Live Board

A real-time multi-room chat application built with Go and vanilla JavaScript.

**Live Demo:** [https://live-board.onrender.com](https://live-board.onrender.com)

## Features

- Real-time messaging via WebSockets
- Multiple room support with create/delete functionality
- User authentication with bcrypt password hashing
- Message persistence (last 200 messages per room)
- Live user count display and presence system
- Auto-created accounts on first login
- Deployment on Render

## Tech Stack

- **Backend:** Go (Gorilla WebSocket)
- **Frontend:** Vanilla JavaScript, HTML, CSS
- **Database:** PostgreSQL
- **Deployment:** Render

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/konradfoulk/live_board.git
   cd live_board
   ```

2. Create the PostgreSQL database:
   ```bash
   psql -U postgres -c "CREATE DATABASE liveboard;"
   ```

3. Set the database connection string and run:
   ```bash
   # Windows (Command Prompt)
   set DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/liveboard?sslmode=disable
   go run server.go models.go db.go

   # Windows (PowerShell)
   $env:DATABASE_URL="postgres://postgres:yourpassword@localhost:5432/liveboard?sslmode=disable"
   go run server.go models.go db.go

   # Linux/macOS
   DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/liveboard?sslmode=disable go run server.go models.go db.go
   ```

4. Open [http://localhost:8080](http://localhost:8080)

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/liveboard?sslmode=disable` |

**Note:** Special characters in passwords must be URL-encoded (e.g., `@` becomes `%40`).

## Architecture

The application uses a Hub pattern for managing WebSocket connections:

- **Hub** - Central broker managing all clients and rooms
- **Client** - Represents a WebSocket connection with dedicated read/write goroutines
- **Room** - Chat room with broadcast channel for message distribution
- All real-time communication is done over **WebSockets**
- All other operations are completed through the **HTTP REST API**

