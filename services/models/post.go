package models

import (
	"so-many-v2/realtime_comments/pkg/errors/agg_errors"
	"time"
)

type Post struct {
	ID       uint      `db:"id" json:"id"`
	AuthorID uint      `db:"user_id" json:"author_id"`
	Title    string    `db:"title" json:"title"`
	Text     string    `db:"text" json:"text"`
	Created  time.Time `db:"created" json:"created"`
	Updated  time.Time `db:"updated" json:"updated"`
}

type CreatePost struct {
	Title  string `db:"title" json:"title"`
	Text   string `db:"text" json:"text"`
	UserID uint   `db:"user_id" json:"user_id"`
}

func (cp *CreatePost) Validate() error {
	if cp.Title == "" {
		return agg_errors.ValidationError{Field: "title", Msg: "empty title"}
	}
	if cp.Text == "" {
		return agg_errors.ValidationError{Field: "text", Msg: "empty text"}
	}
	if cp.UserID == 0 {
		return agg_errors.ValidationError{Field: "user", Msg: "empty user"}
	}
	return nil
}
