package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	General  *GeneralConfig
	Server   *ServerConfig
	Postgres *PostgresConfig
}

type GeneralConfig struct {
	AppName string
	Debug   bool
}

func NewConfig() *Config {
	_ = godotenv.Load()

	debug := false
	debugEnv := os.Getenv("DEBUG")
	if debugEnv == "1" {
		debug = true
	}

	return &Config{
		Server:   NewServerConfig(),
		Postgres: NewPostgresConfig(),
		General: &GeneralConfig{
			AppName: "Comment service",
			Debug:   debug,
		},
	}
}
