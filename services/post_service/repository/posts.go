package repository

import (
	"context"
	"database/sql"
	"errors"
	"so-many-v2/realtime_comments/pkg/errors/db_errors"
	"so-many-v2/realtime_comments/services/models"
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

func (r *PgRepo) ListPosts(ctx context.Context, limit, offset uint) ([]models.Post, error) {
	stmt := "SELECT * FROM posts ORDER BY created DESC LIMIT $1 OFFSET $2"

	posts := make([]models.Post, 0, limit)
	if err := r.DB.SelectContext(ctx, &posts, stmt, limit, offset); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PgRepo) CreatePost(ctx context.Context, post *models.CreatePost) (uint, error) {
	stmt := "INSERT INTO posts (user_id, title, text) VALUES ($1, $2, $3) RETURNING id"

	var postID uint
	err := r.DB.QueryRowContext(ctx, stmt, post.UserID, post.Title, post.Text).Scan(&postID)
	if err != nil {
		return 0, err
	}

	return postID, nil
}
