package db

import (
	"database/sql"
	"embed"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var embedMigrations embed.FS

func getDBPath() string {
	if path := os.Getenv("DB_PATH"); path != "" {
		return path
	}
	return "doppler.db" // fallback to relative path
}

func Connect() *sql.DB {
	dbPath := getDBPath()
	
	// Ensure directory exists
	if dir := os.Getenv("DB_PATH"); dir != "" {
		// Extract directory from path like "data/doppler.db"
		if idx := strings.LastIndex(dir, "/"); idx != -1 {
			dirPath := dir[:idx]
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				log.Panicf("Failed to create directory %v: %v", dirPath, err)
			}
		}
	}
	
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Panicf("Failed to connect to %v: %v", dbPath, err)
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
