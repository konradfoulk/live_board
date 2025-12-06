-- Users table
CREATE TABLE users (
    id
    username
    password_hash
);

-- Rooms table
CREATE TABLE rooms (
    id
    name
    created_at
);

-- Messages table
CREATE TABLE messages (
    id
    room
    user
    content
    created_at
    FOREIGN KEY (room) REFERENCES rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (user) REFERENCES users(id)
);