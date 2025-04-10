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

	// db.MustExec(`
	// 	DROP TABLE IF EXISTS users;
	// `)
	// db.MustExec(`
	// 	CREATE TABLE IF NOT EXISTS users (
	// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 		name TEXT NOT NULL,
	// 		email TEXT NOT NULL UNIQUE,
	// 		password TEXT NOT NULL,
	// 		mfa BOOLEAN NOT NULL DEFAULT false,
	// 		pub TEXT DEFAULT NULL
	// 	);
	// `)

	// db.MustExec(`UPDATE users SET mfa = false WHERE id = 1`)
	// db.MustExec(`DROP TABLE IF EXISTS mfa_sessions;`)
	// db.MustExec(`
	// 	CREATE TABLE IF NOT EXISTS mfa_sessions (
	// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 		user_id INTEGER NOT NULL,
	// 		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	// 		pos_verified INTEGER NOT NULL DEFAULT false,
	// 		neg_verified INTEGER NOT NULL DEFAULT false,
	// 		match BOOLEAN NOT NULL DEFAULT false,
	// 		used BOOLEAN NOT NULL DEFAULT false,
	// 		used_at TIMESTAMP NULL,
	// 		FOREIGN KEY (user_id) REFERENCES users(id)
	// 	);
	// `)

	db.MustExec(`
	CREATE TABLE IF NOT EXISTS register_session (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		used BOOLEAN NOT NULL DEFAULT false,
		used_at TIMESTAMP NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`)
}
