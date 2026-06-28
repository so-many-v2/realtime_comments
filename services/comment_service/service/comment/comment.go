package comment

import (
	"context"
	"encoding/json"
	"fmt"
	"so-many-v2/realtime_comments/pkg/errors/agg_errors"
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/models"
	"time"
)

type CommentsRepoI interface {
	GetCommentButch(ctx context.Context, postID uint, timeFrom time.Time, limit uint) ([]models.Comment, error)
	CreateComment(ctx context.Context, comment *models.CreateComment) (*models.Comment, error)
}

type RedisClientI interface {
	Publish(ctx context.Context, channel string, payload []byte) error
}

type CommentsService struct {
	logger *logg.Logger
	repo   CommentsRepoI
	redis  RedisClientI
}

func NewCommentsService(logger *logg.Logger, repo CommentsRepoI, redis RedisClientI) *CommentsService {
	return &CommentsService{
		logger: logger,
		repo:   repo,
		redis:  redis,
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

	comment, err := cs.repo.CreateComment(ctx, data)
	if err != nil {
		cs.logger.WithField("event", "create_comment").WithError(err).Error("create comments error")
		return 0, err
	}

	_ = cs.sendComment(ctx, comment)
	return comment.ID, err
}

func (cs *CommentsService) sendComment(ctx context.Context, comment *models.Comment) error {
	queueName := fmt.Sprintf("post:%d", comment.ID)
	payload, err := json.Marshal(comment)
	if err != nil {
		cs.logger.WithField("event", "stringify comment").WithError(err).Error("error stringify comment")
		return err
	}

	if err := cs.redis.Publish(ctx, queueName, payload); err != nil {
		cs.logger.WithField("event", "publish comment").WithError(err).Error("error publish comment")
		return err
	}
	return nil
}
