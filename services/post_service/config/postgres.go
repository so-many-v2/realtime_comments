package config

import (
	"fmt"
	"os"
	"time"
)

type PostgresConfig struct {
	user     string
	password string
	host     string
	port     string
	name     string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		user:     os.Getenv("POSTGRES_USER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		host:     os.Getenv("POSTGRES_HOST"),
		port:     os.Getenv("POSTGRES_PORT"),
		name:     os.Getenv("POSTGRES_NAME"),

		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}
}

func (pc *PostgresConfig) Dsn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		pc.user,
		pc.password,
		pc.host,
		pc.port,
		pc.name,
	)
}
