package config

import "os"

type ServerConfig struct {
	Port string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: os.Getenv("POST_SERVICE_PORT"),
	}
}
