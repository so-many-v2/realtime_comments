package test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/models"
	"so-many-v2/realtime_comments/services/post_service/service/posts"
	"so-many-v2/realtime_comments/services/post_service/service/posts/test/mocks"
)

func NewPostServiceTest(repo posts.PostRepoI) *posts.PostService {
	logger := logg.NewLogger("test post service", logg.LevelInfo, false)
	return posts.NewPostService(logger, repo)
}

func TestGetPost_OK(t *testing.T) {
	want := &models.Post{ID: 1, Title: "hello", Text: "body"}

	repo := mocks.NewPostRepoI(t)
	repo.EXPECT().
		GetPostById(mock.Anything, uint(1)).
		Return(want, nil).
		Once()

	got, err := NewPostServiceTest(repo).GetPost(context.Background(), 1)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestGetPost_RepoError(t *testing.T) {
	repoErr := errors.New("db down")

	repo := mocks.NewPostRepoI(t)
	repo.EXPECT().
		GetPostById(mock.Anything, uint(1)).
		Return(nil, repoErr).
		Once()

	_, err := NewPostServiceTest(repo).GetPost(context.Background(), 1)

	require.ErrorIs(t, err, repoErr)
}

func TestCreatePost_OK(t *testing.T) {
	repo := mocks.NewPostRepoI(t)
	repo.EXPECT().
		CreatePost(mock.Anything, mock.Anything).
		Run(func(ctx context.Context, p *models.CreatePost) {
			require.Equal(t, "title", p.Title)
			require.Equal(t, "text", p.Text)
			require.Equal(t, uint(5), p.UserID)
		}).
		Return(uint(7), nil).
		Once()

	id, err := NewPostServiceTest(repo).CreatePost(context.Background(), &models.CreatePost{
		Title:  "title",
		Text:   "text",
		UserID: 5,
	})

	require.NoError(t, err)
	require.Equal(t, uint(7), id)
}

func TestCreatePost_RepoError(t *testing.T) {
	repoErr := errors.New("insert failed")

	repo := mocks.NewPostRepoI(t)
	repo.EXPECT().
		CreatePost(mock.Anything, mock.Anything).
		Return(uint(0), repoErr).
		Once()

	id, err := NewPostServiceTest(repo).CreatePost(context.Background(), &models.CreatePost{
		Title:  "t",
		Text:   "body",
		UserID: 1,
	})

	require.ErrorIs(t, err, repoErr)
	require.Zero(t, id)
}
