package http

import (
	"net/http"
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/comment_service/delivery/http/handlers"
	"so-many-v2/realtime_comments/services/comment_service/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	logger  *logg.Logger
	service *service.Service
}

func NewRouter(logg *logg.Logger, serv *service.Service) *Router {
	return &Router{
		logger:  logg,
		service: serv,
	}
}

func (ro *Router) Init() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	handler := handlers.NewRouterHandler(ro.logger, ro.service)

	router.Group(func(r chi.Router) {
		r.Get("/comments", handler.GetPostHandler)
		r.Post("/comments", handler.CreatePostHandler)
	})

	return router
}
