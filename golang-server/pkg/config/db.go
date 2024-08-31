package config

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func ConnectDB() {
	db, err := sqlx.Open("sqlite3", "file:db.sqlite?cache=shared&mode=rwc&_journal_mode=WAL")
	if err != nil {
		panic(err)
	}
	DB = db
}
