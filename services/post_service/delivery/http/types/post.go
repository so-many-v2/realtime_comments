package types

import (
	"so-many-v2/realtime_comments/pkg/errors/agg_errors"
	"so-many-v2/realtime_comments/services/post_service/repository/models"
)

type CreatePostSchema struct {
	Text     string `json:"text"`
	Title    string `json:"title"`
	AuthorID uint   `json:"author_id"` // TODO: Send as JWT
}

func (cp *CreatePostSchema) Validate() error {
	if cp.Title == "" {
		return agg_errors.ValidationError{Field: "title", Msg: "empty title"}
	}
	if cp.Text == "" {
		return agg_errors.ValidationError{Field: "text", Msg: "empty text"}
	}
	if cp.AuthorID == 0 {
		return agg_errors.ValidationError{Field: "author", Msg: "empty author"}
	}
	return nil
}

func (ps *CreatePostSchema) ToModel() *models.CreatePost {
	return &models.CreatePost{
		Title:    ps.Title,
		Text:     ps.Text,
		AuthorID: ps.AuthorID,
	}
}
