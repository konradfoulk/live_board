-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL
);

-- Rooms table
CREATE TABLE rooms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Messages table
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_type TEXT NOT NULL, -- "join", "leave", "chat"
    room INTEGER NOT NULL,
    user INTEGER NOT NULL,
    content TEXT, -- NULL for join and leave messages
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (room) REFERENCES rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (user) REFERENCES users(id)
);