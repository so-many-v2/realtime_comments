package service

import (
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/comment_service/clients/redis_client"
	"so-many-v2/realtime_comments/services/comment_service/repository/pg_repo"
	"so-many-v2/realtime_comments/services/comment_service/service/comment"
)

type Service struct {
	Comments *comment.CommentsService
}

func NewService(logg *logg.Logger, repo *pg_repo.PgRepo, redis *redis_client.RedisClient) *Service {
	return &Service{
		Comments: comment.NewCommentsService(logg, repo, redis),
	}
}
