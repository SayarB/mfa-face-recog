package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}
	fmt.Println(os.Getenv("CGO_ENABLED"))
	db, err := sqlx.Open("sqlite3", "file:db.sqlite?cache=shared&mode=rwc&_journal_mode=WAL")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// db.MustExec(`
	// 	CREATE TABLE IF NOT EXISTS users (
	// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 		name TEXT NOT NULL,
	// 		email TEXT NOT NULL UNIQUE,
	// 		password TEXT NOT NULL
	// 	);
	// `)

	db.MustExec(`
		DELETE FROM users;
	`)
}
