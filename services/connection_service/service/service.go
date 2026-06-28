package service

import (
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/connection_service/service/connection"
)

type Service struct {
	Connections *connection.ConnectionService
}

func NewService(logger *logg.Logger) *Service {
	return &Service{
		Connections: connection.NewConnectionService(logger),
	}
}
