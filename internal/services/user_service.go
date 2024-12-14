package services

import (
	"database/sql"
	"doppler/internal/models"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func GetUsers(db *sql.DB) []models.User {

	var query, err = db.Prepare("SELECT u.id, u.username, u.email, u.created FROM user u;")

	if err != nil {
		log.Panic(err)
	}

	rows, err := query.Query()
	if err != nil {
		log.Panicf("Query failed: %v", err)
	}
	defer rows.Close()
	users := []models.User{}
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.Created)
		if err != nil {
			log.Panicf("Failed scanning to Post: %v", err)
		}
		users = append(users, u)
	}
	return users
}

func CreateUser(db *sql.DB, username string, password string, email string) models.User {
	var query, err = db.Prepare("INSERT INTO user (username, password, email) VALUES (?, ?, ?) RETURNING id, username, email")
	if err != nil {
		log.Panic(err)
	}
	cp := &cryptoParams{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
	user := models.User{}
	passwordHash, err := generateFromPassword(password, cp)
	if err != nil {
		log.Panicf("Failed generating password: %v", err)
	}
	err = query.QueryRow(username, passwordHash, email).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		log.Panicf("Query failed: %v", err)
	}
	return user

}

func ValidateUser(db *sql.DB, username string, password string) bool {

	var query, err = db.Prepare("SELECT u.password, FROM user u WHERE u.username = ?;")

	if err != nil {
		log.Panic(err)
	}
	var encodedHash string
	err = query.QueryRow(username).Scan(&encodedHash)
	if err != nil {
		log.Panicf("Query failed: %v", err)
	}
	comparePasswordAndHash(password, encodedHash)
	return true
}
