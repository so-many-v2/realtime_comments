package http_tools

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type HttpServer struct {
	server *http.Server
}

func NewHttpServer(address string, router http.Handler) *HttpServer {
	return &HttpServer{
		server: &http.Server{
			Addr:           address,
			Handler:        router,
			MaxHeaderBytes: MaxHeadersSize,
			ReadTimeout:    time.Second * ReadTimeOutSec,
			WriteTimeout:   time.Second * WriteTimeOutSec,
		},
	}
}

func NewStreamingServer(address string, router http.Handler) *HttpServer {
	return &HttpServer{
		server: &http.Server{
			Addr:           address,
			Handler:        router,
			MaxHeaderBytes: MaxHeadersSize,
			ReadTimeout:    time.Second * ReadTimeOutSec,
			WriteTimeout:   0,
			IdleTimeout:    0,
		},
	}
}

func (hs *HttpServer) StartServer() error {
	if err := hs.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (hs *HttpServer) Shutdown(ctx context.Context) error {
	if err := hs.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
