package redis_client

import (
	"context"
	"fmt"
	"time"

	"so-many-v2/realtime_comments/services/comment_service/config"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(ctx context.Context, cfg *config.RedisConfig) (*RedisClient, error) {
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

	return &RedisClient{Client: client}, nil
}

func (rc *RedisClient) Publish(ctx context.Context, channel string, payload []byte) error {
	if err := rc.Client.Publish(ctx, channel, payload).Err(); err != nil {
		return fmt.Errorf("redis: publish to %q: %w", channel, err)
	}
	return nil
}
