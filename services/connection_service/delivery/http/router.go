package http

import (
	"net/http"

	"so-many-v2/realtime_comments/pkg/http_tools"
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/connection_service/delivery/http/handlers"
	"so-many-v2/realtime_comments/services/connection_service/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	logger  *logg.Logger
	handler *handlers.RouterHandler
}

func NewRouter(logger *logg.Logger, serv *service.Service) *Router {
	return &Router{
		logger:  logger,
		handler: handlers.NewRouterHandler(logger, serv),
	}
}

func (ro *Router) Init() http.Handler {
	router := chi.NewRouter()

	router.Use(http_tools.CORS)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	router.Get("/api/sse/{post_id}", ro.handler.SSEHandler)

	return router
}
