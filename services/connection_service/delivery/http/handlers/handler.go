package handlers

import (
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/connection_service/service"
)

type RouterHandler struct {
	logger  *logg.Logger
	service *service.Service
}

func NewRouterHandler(logger *logg.Logger, serv *service.Service) *RouterHandler {
	return &RouterHandler{
		logger:  logger,
		service: serv,
	}
}