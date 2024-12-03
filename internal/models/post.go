package models

import "time"

type Post struct {
	ID      int       `json:"id"`
	Content string    `json:"content"`
	User    string    `json:"user"`
	Created time.Time `json:"created"`
}
