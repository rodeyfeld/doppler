package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var dbName string = "doppler.db"

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Panicf("Failed to connect to %v: %v", dbName, err)
	}
	return db
}
