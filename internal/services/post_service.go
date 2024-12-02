package services

import (
	"database/sql"
	"doppler/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func GetPosts(db *sql.DB) []models.Post {

	var query, err = db.Prepare("SELECT p.ID, p.NAME FROM post p;")

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
		err = rows.Scan(&p.ID, &p.Name)
		if err != nil {
			log.Panicf("Failed scanning to Post: %v", err)
		}
		posts = append(posts, p)
	}
	return posts
}

func CreatePost(db *sql.DB, name string) models.Post {

	var query, err = db.Prepare("INSERT INTO post (name) VALUES (?) RETURNING *")
	if err != nil {
		log.Panic(err)
	}

	post := models.Post{}
	err = query.QueryRow(name).Scan(&post.ID, &post.Name)
	if err != nil {
		log.Panicf("Query failed: %v", err)
	}

	return post
}
