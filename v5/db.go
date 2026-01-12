package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func initDatabase() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/liveboard?sslmode=disable"
	}

	db, _ := sql.Open("postgres", connStr)

	schema, _ := os.ReadFile("schema.sql")
	db.Exec(string(schema))

	return db
}
