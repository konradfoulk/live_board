package main

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func initDatabase() *sql.DB {
	db, _ := sql.Open("sqlite", "./app.db")

	schema, _ := os.ReadFile("schema.sql")
	db.Exec(string(schema))

	return db
}
