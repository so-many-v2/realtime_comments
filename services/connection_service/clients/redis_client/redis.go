package redis_client

import (
	"context"
	"fmt"
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/connection_service/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	logger *logg.Logger
	*redis.Client
}

func NewRedisClient(ctx context.Context, logger *logg.Logger, cfg config.RedisConfig) (*RedisClient, error) {
	opt, err := redis.ParseURL(cfg.Dsn())
	if err != nil {
		return nil, fmt.Errorf("redis: parse url: %w", err)
	}

	client := redis.NewClient(opt)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis: ping: %w", err)
	}

	return &RedisClient{
		Client: client,
		logger: logger,
	}, nil
}


func (rc *RedisClient) SubscribeChannel(ctx context.Context, pattern string) (<-chan *redis.Message, error) {
	sub := rc.Client.PSubscribe(ctx, pattern)

	if _, err := sub.Receive(ctx); err != nil {
		rc.logger.WithField("event", "psubscribe").
			WithError(err).
			Error(fmt.Sprintf("error psubscribe to %s pattern", pattern))
		return nil, err
	}

	return sub.Channel(), nil
}
