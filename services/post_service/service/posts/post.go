package posts

import (
	"context"
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/post_service/delivery/http/types"
	"so-many-v2/realtime_comments/services/post_service/repository/models"
)

//go:generate mockery --name=PostRepoI --output=./test/mocks --outpkg=mocks --filename=PostRepoI.go --with-expecter
type PostRepoI interface {
	GetPostById(ctx context.Context, postID uint) (*models.Post, error)
	CreatePost(ctx context.Context, post *models.CreatePost) (uint, error)
}
type PostService struct {
	logger *logg.Logger
	repo   PostRepoI
}

func NewPostService(logger *logg.Logger, repo PostRepoI) *PostService {
	return &PostService{
		logger: logger,
		repo:   repo,
	}
}

func (ps *PostService) CreatePost(ctx context.Context, data *types.CreatePostSchema) (uint, error) {
	if err := data.Validate(); err != nil {
		return 0, err
	}

	model := data.ToModel()
	postId, err := ps.repo.CreatePost(ctx, model)
	if err != nil {
		ps.logger.WithError(err).Error("create post failed")
		return 0, err
	}
	return postId, nil
}

func (ps *PostService) GetPost(ctx context.Context, postID uint) (*models.Post, error) {
	post, err := ps.repo.GetPostById(ctx, postID)
	if err != nil {
		return nil, err
	}
	return post, err
}
