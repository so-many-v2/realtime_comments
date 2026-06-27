package repository

import (
	"context"
	"so-many-v2/realtime_comments/services/models"
	"time"
)

func (r *PgRepo) GetCommentButch(ctx context.Context, postID uint, timeFrom time.Time, limit uint) ([]models.Comment, error) {
	stmt := `SELECT * FROM comments c
			 JOIN post_commetns pc on pc.comment_id = c.id
			 WHERE pc.post_id = $1 AND c.created > $2
			 ORDER BY c.created 
			 LIMIT $3`

	comments := make([]models.Comment, 0, limit)
	if err := r.DB.SelectContext(ctx, comments, stmt, postID, timeFrom, limit); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *PgRepo) CreateComment(ctx context.Context, comment *models.CreateComment) (uint, error) {
	trx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer trx.Rollback()

	insertCommentStmt := "INSERT INTO comment (text, user_id) VALUES ($1, $2) RETURNING id"
	insertM2mStmt := "INSERT INTO post_comments (post_id, comment_id) VALUES ($1, $2) RETURNING id"

	var commentID uint
	if err = trx.QueryRowContext(
		ctx, insertCommentStmt,
		comment.Text,
		comment.UserID,
	).Scan(&commentID); err != nil {
		return 0, err
	}

	if _, err = trx.ExecContext(
		ctx,
		insertM2mStmt,
		comment.PostID,
		commentID,
	); err != nil {
		return 0, err
	}

	if err = trx.Commit(); err != nil {
		return 0, err
	}

	return commentID, nil
}
