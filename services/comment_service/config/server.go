package config

import "os"

type ServerConfig struct {
	Port string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: os.Getenv("COMMENT_SERVICE_PORT"),
	}
}
