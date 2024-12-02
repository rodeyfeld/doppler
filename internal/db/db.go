package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var dbName string = "doppler.db"

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Panicf("Failed to connect to %v: %v", dbName, err)
	}
	return db
}

func SetupDb() {
	if _, err := os.Stat(dbName); errors.Is(err, os.ErrNotExist) {
		Initialize()
	}
}

func Initialize() {
	file, err := os.Create(dbName)
	defer file.Close()
	if err != nil {
		log.Panicf("Failed to create file %v: %v", dbName, err)
	}
	log.Printf("Created file %s", dbName)
}
