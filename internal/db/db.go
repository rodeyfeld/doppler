package db

import (
	"database/sql"
	"embed"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var embedMigrations embed.FS

var dbName string = "doppler.db"

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Panicf("Failed to connect to %v: %v", dbName, err)
	}
	return db
}

// RunMigrations applies all pending migrations to the database
func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(db, "sql"); err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
