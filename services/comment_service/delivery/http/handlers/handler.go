package handlers

import (
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/comment_service/service"
)

type RouterHandler struct {
	logger  *logg.Logger
	service *service.Service
}

func NewRouterHandler(logg *logg.Logger, serv *service.Service) *RouterHandler {
	return &RouterHandler{
		logger:  logg,
		service: serv,
	}
}
