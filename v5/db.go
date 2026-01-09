package main

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func initDatabase() *sql.DB {
	db, _ := sql.Open("sqlite3", "./app.db")

	schema, _ := os.ReadFile("schema.sql")
	db.Exec(string(schema))

	return db
}
