# TODO.md

## Project Goal

Refactor the application to use PostgreSQL instead of SQLite for the database.

## Current State

- Database driver: `modernc.org/sqlite` (pure Go SQLite)
- Database file: `./app.db`
- PostgreSQL driver already in go.mod: `github.com/lib/pq v1.10.9`

## Files Requiring Changes

### db.go
- Change driver import from `modernc.org/sqlite` to `github.com/lib/pq`
- Update `sql.Open()` to use PostgreSQL connection string
- Add environment variable or config for connection string

### schema.sql
SQLite-specific syntax that needs PostgreSQL equivalents:
- `INTEGER PRIMARY KEY AUTOINCREMENT` → `SERIAL PRIMARY KEY`
- `DATETIME DEFAULT CURRENT_TIMESTAMP` → `TIMESTAMP DEFAULT CURRENT_TIMESTAMP`
- `TEXT` works in both but consider `VARCHAR` for constrained fields

### server.go
Placeholder syntax changes (`?` → `$1, $2, ...`):
- Line 44: `INSERT OR IGNORE INTO rooms` → `INSERT INTO rooms ... ON CONFLICT DO NOTHING`
- Line 97: `INSERT INTO rooms (name) VALUES (?)`
- Line 140: `DELETE FROM rooms WHERE name = ?`
- Line 193: `SELECT id, password_hash FROM users WHERE username = ?`
- Line 200: `INSERT INTO users (username, password_hash) VALUES (?, ?)`
- Line 237: `SELECT name FROM rooms ORDER BY created_at ASC` (no change needed)

### models.go
Placeholder syntax changes:
- Line 123: `SELECT` query for messages (check for `?` placeholders)
- Line 158: `INSERT INTO messages` with `?` placeholders

## PostgreSQL Setup Requirements

- PostgreSQL server running (local or remote)
- Database created for the application
- Connection string format: `postgres://user:password@host:port/dbname?sslmode=disable`

## Migration Considerations

- Existing data in `app.db` will need manual migration if preserving data
- Schema should use `IF NOT EXISTS` for idempotent startup
