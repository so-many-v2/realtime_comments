package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	General *GeneralConfig
	Server  *ServerConfig
	Redis   *RedisConfig
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
		Server: NewServerConfig(),
		Redis:  NewRedisConfig(),
		General: &GeneralConfig{
			AppName: "Connection service",
			Debug:   debug,
		},
	}
}
