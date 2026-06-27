package repository

import (
	"context"
	"database/sql"
	"errors"
	"so-many-v2/realtime_comments/pkg/errors/db_errors"
	"so-many-v2/realtime_comments/services/post_service/repository/models"
)

func (r *PgRepo) GetPostById(ctx context.Context, postID uint) (*models.Post, error) {
	stmt := "SELECT * FROM posts WHERE id = $1"

	var post models.Post
	if err := r.DB.GetContext(ctx, &post, stmt, postID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, db_errors.ErrNotFound
		}
		return nil, err
	}

	return &post, nil
}

func (r *PgRepo) CreatePost(ctx context.Context, post *models.CreatePost) (uint, error) {
	stmt := "INSERT INTO posts (title, text) VALUES ($1, $2) RETURNING id"

	var postID uint
	err := r.DB.QueryRowContext(ctx, stmt, post.Title, post.Text).Scan(&postID)
	if err != nil {
		return 0, err
	}

	return postID, nil
}
