package services

import (
	"database/sql"
	"doppler/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func GetPosts(db *sql.DB) []models.Post {

	var query, err = db.Prepare("SELECT p.user_id, p.title, p.content, p.created FROM post p;")

	if err != nil {
		log.Panic(err)
	}

	rows, err := query.Query()
	if err != nil {
		log.Panicf("Query failed: %v", err)
	}
	defer rows.Close()
	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		p.UserID = 1
		err = rows.Scan(&p.UserID, &p.Title, &p.Content, &p.Created)
		if err != nil {
			log.Panicf("Failed scanning to Post: %v", err)
		}
		posts = append(posts, p)
	}
	return posts
}

func CreatePost(db *sql.DB, userID int, title string, content string) models.Post {

	var query, err = db.Prepare("INSERT INTO post (user_id, title, content) VALUES (?, ?, ?) RETURNING id, user_id, title, content")
	if err != nil {
		log.Panic(err)
	}

	post := models.Post{}
	err = query.QueryRow(userID, title, content).Scan(&post.ID, &post.UserID, &post.Title, &post.Content)
	if err != nil {
		log.Panicf("Query failed: %v", err)
	}

	return post
}
