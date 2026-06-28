package pg_repo

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

func (r *PgRepo) CreateComment(ctx context.Context, comment *models.CreateComment) (*models.Comment, error) {
	trx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer trx.Rollback()

	insertCommentStmt := "INSERT INTO comment (text, user_id) VALUES ($1, $2) RETURNING *"
	insertM2mStmt := "INSERT INTO post_comments (post_id, comment_id) VALUES ($1, $2) RETURNING id"

	var newComment models.Comment
	if err = trx.QueryRowContext(
		ctx, insertCommentStmt,
		comment.Text,
		comment.UserID,
	).Scan(&comment); err != nil {
		return nil, err
	}

	if _, err = trx.ExecContext(
		ctx,
		insertM2mStmt,
		comment.PostID,
		newComment.ID,
	); err != nil {
		return nil, err
	}

	if err = trx.Commit(); err != nil {
		return nil, err
	}

	return &newComment, nil
}
