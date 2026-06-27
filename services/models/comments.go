package models

import (
	"so-many-v2/realtime_comments/pkg/errors/agg_errors"
	"time"
)

type Comment struct {
	ID      uint      `db:"id" json:"id"`
	UserID  uint      `db:"user_id" json:"author_id"`
	PostID  uint      `db:"post_id" json:"post_id"`
	Text    string    `db:"text" json:"text"`
	Created time.Time `db:"created" json:"created"`
	Updated time.Time `db:"updated" json:"updated"`
}

type CreateComment struct {
	Text   string `db:"text" json:"text"`
	UserID uint   `db:"user_id" json:"user_id"`
	PostID uint   `db:"post_id" json:"post_id"`
}

func (cp *CreateComment) Validate() error {
	if cp.Text == "" {
		return agg_errors.ValidationError{Field: "text", Msg: "empty text"}
	}
	if cp.PostID == 0 {
		return agg_errors.ValidationError{Field: "post_id", Msg: "empty post_id"}
	}
	if cp.UserID == 0 {
		return agg_errors.ValidationError{Field: "user", Msg: "empty user"}
	}
	return nil
}
