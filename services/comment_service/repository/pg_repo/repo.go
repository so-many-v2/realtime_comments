package pg_repo

import (
	"context"
	"fmt"
	"time"

	"so-many-v2/realtime_comments/services/comment_service/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PgRepo struct {
	DB *sqlx.DB
}

func NewPostgres(ctx context.Context, cfg *config.PostgresConfig) (*PgRepo, error) {
	db, err := sqlx.Open("pgx", cfg.Dsn())
	if err != nil {
		return nil, fmt.Errorf("postgres: open: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}

	return &PgRepo{DB: db}, nil
}
