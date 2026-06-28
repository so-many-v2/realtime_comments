package http

import (
	"net/http"
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/post_service/delivery/http/handlers"
	"so-many-v2/realtime_comments/services/post_service/service"

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

	router.Route("/api", func(r chi.Router) {
		r.Get("/posts", handler.ListPostsHandler)
		r.Get("/posts/{postID}", handler.GetPostHandler)
		r.Post("/posts", handler.CreatePostHandler)
	})

	return router
}
