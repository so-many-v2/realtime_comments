package models

import (
	"time"
)

type Post struct {
	ID       uint      `db:"id" json:"id"`
	AuthorID uint      `db:"author_id" json:"author_id"`
	Title    string    `db:"title" json:"title"`
	Text     string    `db:"text" json:"text"`
	Created  time.Time `db:"created" json:"created"`
	Updated  time.Time `db:"updated" json:"updated"`
}

type CreatePost struct {
	Title    string `db:"title"`
	Text     string `db:"text"`
	AuthorID uint   `db:"author_id"`
}
