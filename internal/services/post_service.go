package services

import (
	"database/sql"
	"doppler/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func GetPosts(db *sql.DB) []models.Post {

	var query, err = db.Prepare("SELECT p.user, p.content, p.created FROM post p;")

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
		p.User = "User0"
		err = rows.Scan(&p.User, &p.Content, &p.Created)
		if err != nil {
			log.Panicf("Failed scanning to Post: %v", err)
		}
		posts = append(posts, p)
	}
	return posts
}

func CreatePost(db *sql.DB, user string, content string) models.Post {

	var query, err = db.Prepare("INSERT INTO post (user, content) VALUES (?, ?) RETURNING id, content, user")
	if err != nil {
		log.Panic(err)
	}

	post := models.Post{}
	err = query.QueryRow(user, content).Scan(&post.ID, &post.Content, &post.User)
	if err != nil {
		log.Panicf("Query failed: %v", err)
	}

	return post
}
