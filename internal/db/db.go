package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type DB struct {
	Db *sql.DB
}

var dbName string = "doppler.db"

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to %s: %s", dbName, err)
	}
	return db, nil
}

func SetupDb() {
	if _, err := os.Stat(dbName); errors.Is(err, os.ErrNotExist) {
		Initialize()
		Migrate()
	}
}

func Initialize() error {
	file, err := os.Create(dbName)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("Failed to create file %s: %s", dbName, err)
	}
	log.Printf("Created file %s", dbName)
	return nil
}
func handleMigration(conn *sql.DB, path string, info fs.FileInfo, err error) error {
	if err != nil {
		return err
	}
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to read migration sql file %s: %s", path, err)
	}
	_, err = conn.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("Failed to execute migration sql file %s: %s", path, err)
	}
	return nil
}
func Migrate() {
	conn, err := Connect()
	if err != nil {
		log.Panicf("Could not connect to DB during migration! %s", err)
	}
	filepath.Walk("sql/migrations", func(path string, entry fs.FileInfo, err error) error {
		return handleMigration(conn, path, entry, err)
	})
}
