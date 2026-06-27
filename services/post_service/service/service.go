package service

import (
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/post_service/repository"
	"so-many-v2/realtime_comments/services/post_service/service/posts"
)

type Service struct {
	Post *posts.PostService
}

func NewService(logg *logg.Logger, repo *repository.PgRepo) *Service {
	return &Service{
		Post: posts.NewPostService(logg, repo),
	}
}
