package repository

import (
	"context"
	"so-many-v2/realtime_comments/services/comment_service/config"
	"so-many-v2/realtime_comments/services/comment_service/repository/db"
)

type PgRepo struct {
	*db.Postgres
}

func NewPgRepo(ctx context.Context, cfg *config.PostgresConfig) (*PgRepo, error) {
	postgres, err := db.NewPostgres(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &PgRepo{
		postgres,
	}, nil
}
