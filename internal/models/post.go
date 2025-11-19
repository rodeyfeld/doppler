package models

import (
	"time"
)

type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	UserID      int       `json:"user_id"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
	PictureURLs []string  `json:"picture_urls,omitempty"` // Presigned URLs for all pictures
}

type Picture struct {
	ID       int    `json:"id"`
	PostID   int    `json:"post_id"`
	Filename string `json:"filename"`
}
