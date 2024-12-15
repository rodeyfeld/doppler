package models

import (
	"time"
)

type Post struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	UserID   int       `json:"user_id"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}
