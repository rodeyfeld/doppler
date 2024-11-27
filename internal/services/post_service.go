package services

import (
	"database/sql"
	"doppler/internal/models"
)

func GetPosts(db *sql.DB) []Post {

	data, err := os.ReadFile("sql/post/get_posts.sql")
	if err != nil {
		log.Panic(err)
	}
	rows, err := db.Query(string(data))
	if err != nil {
		log.Panicf("Query failed: %v", err)
	}
	defer rows.Close()
	posts := []Post{}
	for rows.Next() {
		var p Post
		err = rows.Scan(&p.ID, &p.Name)
		if err != nil {
			log.Panicf("Failed scanning to Post: %v", err)
		}
		posts = append(posts, p)
	}
	return posts
}
