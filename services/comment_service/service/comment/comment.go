package comment

import (
	"context"
	"so-many-v2/realtime_comments/pkg/errors/agg_errors"
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/models"
	"time"
)

type CommentsRepoI interface {
	GetCommentButch(ctx context.Context, postID uint, timeFrom time.Time, limit uint) ([]models.Comment, error)
	CreateComment(ctx context.Context, comment *models.CreateComment) (uint, error)
}

type CommentsService struct {
	logger *logg.Logger
	repo   CommentsRepoI
}

func NewCommentsService(logger *logg.Logger, repo CommentsRepoI) *CommentsService {
	return &CommentsService{
		logger: logger,
		repo:   repo,
	}
}

func (cs *CommentsService) GetPostCommentBatch(
	ctx context.Context,
	postID uint,
	timeFrom time.Time,
	limit uint,
) ([]models.Comment, error) {
	if timeFrom.After(time.Now()) {
		return nil, agg_errors.ValidationError{Field: "time_from", Msg: "time_from can't be in future"}
	}

	comments, err := cs.repo.GetCommentButch(ctx, postID, timeFrom, limit)
	if err != nil {
		cs.logger.WithField("event", "get_comments").WithError(err).Error("get comments error")
		return nil, err
	}
	return comments, err
}

func (cs *CommentsService) CreateComment(ctx context.Context, data *models.CreateComment) (uint, error) {
	if err := data.Validate(); err != nil {
		return 0, err
	}

	commentId, err := cs.repo.CreateComment(ctx, data)
	if err != nil {
		cs.logger.WithField("event", "create_comment").WithError(err).Error("create comments error")
		return 0, err
	}
	return commentId, err
}
