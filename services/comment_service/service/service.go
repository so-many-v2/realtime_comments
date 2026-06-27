package service

import (
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/comment_service/repository"
	"so-many-v2/realtime_comments/services/comment_service/service/comment"
)

type Service struct {
	Comments *comment.CommentsService
}

func NewService(logg *logg.Logger, repo *repository.PgRepo) *Service {
	return &Service{
		Comments: comment.NewCommentsService(logg, repo),
	}
}
