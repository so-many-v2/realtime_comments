package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"so-many-v2/realtime_comments/pkg/http_tools"
	"so-many-v2/realtime_comments/pkg/logg"
	"so-many-v2/realtime_comments/services/post_service/config"
	http_router "so-many-v2/realtime_comments/services/post_service/delivery/http"
	"so-many-v2/realtime_comments/services/post_service/repository"
	"so-many-v2/realtime_comments/services/post_service/service"
	"syscall"
	"time"
)

func StartApp() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.NewConfig()
	logger := logg.NewLogger(
		cfg.General.AppName,
		logg.LevelInfo,
		cfg.General.Debug,
	)

	repo, err := repository.NewPgRepo(ctx, cfg.Postgres)
	if err != nil {
		logger.WithError(err).Fatal("postgres init failed")
	}
	serv := service.NewService(logger, repo)
	router := http_router.NewRouter(logger, serv)
	httpServer := http_tools.NewHttpServer(
		fmt.Sprintf(":%s", cfg.Server.Port),
		router.Init(),
	)

	go func() {
		if err := httpServer.StartServer(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Fatal("http server failed")
		}
	}()

	logger.Infof("App %s started on %s port", cfg.General.AppName, cfg.Server.Port)

	<-shutdown
	logger.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = httpServer.Shutdown(shutdownCtx)
	_ = repo.DB.Close()
}
