package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"so-many-v2/realtime_comments/pkg/errors/agg_errors"
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/comment_service/service/comment"
	"so-many-v2/realtime_comments/services/comment_service/service/comment/test/mocks"
	"so-many-v2/realtime_comments/services/models"
)

func NewCommentsServiceTest(repo comment.CommentsRepoI, redis comment.RedisClientI) *comment.CommentsService {
	logger := logg.NewLogger("test comment service", logg.LevelInfo, false)
	return comment.NewCommentsService(logger, repo, redis)
}

func TestGetPostCommentBatch_OK(t *testing.T) {
	want := []models.Comment{
		{ID: 1, PostID: 7, UserID: 2, Text: "hi"},
		{ID: 2, PostID: 7, UserID: 3, Text: "yo"},
	}

	repo := mocks.NewCommentsRepoI(t)
	repo.EXPECT().
		GetCommentButch(mock.Anything, uint(7), mock.Anything, uint(10)).
		Return(want, nil).
		Once()

	redis := mocks.NewRedisClientI(t)

	got, err := NewCommentsServiceTest(repo, redis).
		GetPostCommentBatch(context.Background(), 7, time.Unix(0, 0), 10)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestGetPostCommentBatch_RepoError(t *testing.T) {
	repoErr := errors.New("db down")

	repo := mocks.NewCommentsRepoI(t)
	repo.EXPECT().
		GetCommentButch(mock.Anything, uint(7), mock.Anything, uint(10)).
		Return(nil, repoErr).
		Once()

	redis := mocks.NewRedisClientI(t)

	_, err := NewCommentsServiceTest(repo, redis).
		GetPostCommentBatch(context.Background(), 7, time.Unix(0, 0), 10)

	require.ErrorIs(t, err, repoErr)
}

func TestGetPostCommentBatch_FutureTime(t *testing.T) {
	repo := mocks.NewCommentsRepoI(t)
	redis := mocks.NewRedisClientI(t)

	_, err := NewCommentsServiceTest(repo, redis).
		GetPostCommentBatch(context.Background(), 7, time.Now().Add(time.Hour), 10)

	var vErr agg_errors.ValidationError
	require.ErrorAs(t, err, &vErr)
	require.Equal(t, "time_from", vErr.Field)
}

func TestCreateComment_OK(t *testing.T) {
	created := &models.Comment{ID: 42, PostID: 7, UserID: 5, Text: "text"}

	repo := mocks.NewCommentsRepoI(t)
	repo.EXPECT().
		CreateComment(mock.Anything, mock.Anything).
		Run(func(ctx context.Context, c *models.CreateComment) {
			require.Equal(t, "text", c.Text)
			require.Equal(t, uint(7), c.PostID)
			require.Equal(t, uint(5), c.UserID)
		}).
		Return(created, nil).
		Once()

	redis := mocks.NewRedisClientI(t)
	redis.EXPECT().
		Publish(mock.Anything, "post:7", mock.Anything).
		Return(nil).
		Once()

	id, err := NewCommentsServiceTest(repo, redis).CreateComment(context.Background(), &models.CreateComment{
		Text:   "text",
		PostID: 7,
		UserID: 5,
	})

	require.NoError(t, err)
	require.Equal(t, uint(42), id)
}

func TestCreateComment_RepoError(t *testing.T) {
	repoErr := errors.New("insert failed")

	repo := mocks.NewCommentsRepoI(t)
	repo.EXPECT().
		CreateComment(mock.Anything, mock.Anything).
		Return(nil, repoErr).
		Once()

	redis := mocks.NewRedisClientI(t)

	id, err := NewCommentsServiceTest(repo, redis).CreateComment(context.Background(), &models.CreateComment{
		Text:   "text",
		PostID: 7,
		UserID: 5,
	})

	require.ErrorIs(t, err, repoErr)
	require.Zero(t, id)
}

func TestCreateComment_ValidationError(t *testing.T) {
	repo := mocks.NewCommentsRepoI(t)
	redis := mocks.NewRedisClientI(t)

	id, err := NewCommentsServiceTest(repo, redis).CreateComment(context.Background(), &models.CreateComment{
		Text:   "",
		PostID: 7,
		UserID: 5,
	})

	var vErr agg_errors.ValidationError
	require.ErrorAs(t, err, &vErr)
	require.Equal(t, "text", vErr.Field)
	require.Zero(t, id)
}

func TestCreateComment_PublishError_StillSucceeds(t *testing.T) {
	created := &models.Comment{ID: 42, PostID: 7, UserID: 5, Text: "text"}

	repo := mocks.NewCommentsRepoI(t)
	repo.EXPECT().
		CreateComment(mock.Anything, mock.Anything).
		Return(created, nil).
		Once()

	redis := mocks.NewRedisClientI(t)
	redis.EXPECT().
		Publish(mock.Anything, "post:7", mock.Anything).
		Return(errors.New("redis down")).
		Once()

	id, err := NewCommentsServiceTest(repo, redis).CreateComment(context.Background(), &models.CreateComment{
		Text:   "text",
		PostID: 7,
		UserID: 5,
	})

	require.NoError(t, err)
	require.Equal(t, uint(42), id)
}
