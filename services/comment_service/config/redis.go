package config

import (
	"fmt"
	"os"
)

type RedisConfig struct {
	user     string
	password string
	host     string
	port     string
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		user:     os.Getenv("REDIS_USER"),
		password: os.Getenv("REDIS_PASSWORD"),
		host:     os.Getenv("REDIS_HOST"),
		port:     os.Getenv("REDIS_PORT"),
	}
}

func (rc *RedisConfig) Dsn() string {
	return fmt.Sprintf("redis://%s:%s@%s:%s/0", rc.user, rc.password, rc.host, rc.port)
}
