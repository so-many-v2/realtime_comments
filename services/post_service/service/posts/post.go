package posts

import (
	"context"
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/models"
)

//go:generate mockery --name=PostRepoI --output=./test/mocks --outpkg=mocks --filename=PostRepoI.go --with-expecter
type PostRepoI interface {
	GetPostById(ctx context.Context, postID uint) (*models.Post, error)
	ListPosts(ctx context.Context, limit, offset uint) ([]models.Post, error)
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

func (ps *PostService) CreatePost(ctx context.Context, data *models.CreatePost) (uint, error) {
	if err := data.Validate(); err != nil {
		return 0, err
	}

	postId, err := ps.repo.CreatePost(ctx, data)
	if err != nil {
		ps.logger.WithField("event", "create_post").WithError(err).Error("create post error")
		return 0, err
	}
	return postId, nil
}

func (ps *PostService) ListPosts(ctx context.Context, limit, offset uint) ([]models.Post, error) {
	const defaultLimit, maxLimit = 20, 100
	if limit == 0 || limit > maxLimit {
		limit = defaultLimit
	}

	posts, err := ps.repo.ListPosts(ctx, limit, offset)
	if err != nil {
		ps.logger.WithField("event", "list_posts").WithError(err).Error("list posts error")
		return nil, err
	}
	return posts, nil
}

func (ps *PostService) GetPost(ctx context.Context, postID uint) (*models.Post, error) {
	post, err := ps.repo.GetPostById(ctx, postID)
	if err != nil {
		ps.logger.WithField("event", "get_post").WithError(err).Error("get post error")
		return nil, err
	}
	return post, err
}
