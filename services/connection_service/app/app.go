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
	"so-many-v2/realtime_comments/services/connection_service/clients/redis_client"
	"so-many-v2/realtime_comments/services/connection_service/config"
	http_router "so-many-v2/realtime_comments/services/connection_service/delivery/http"
	"so-many-v2/realtime_comments/services/connection_service/service"
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

	redisClient, err := redis_client.NewRedisClient(ctx, logger, *cfg.Redis)
	if err != nil {
		logger.WithError(err).Fatal("redis init failed")
	}

	serv := service.NewService(logger)

	messages, err := redisClient.SubscribeChannel(ctx, "post:*")
	if err != nil {
		logger.WithError(err).Fatal("redis subscribe failed")
	}
	go func() {
		for msg := range messages {
			serv.Connections.Broadcast(msg.Channel, []byte(msg.Payload))
		}
	}()

	router := http_router.NewRouter(logger, serv)
	httpServer := http_tools.NewStreamingServer(
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

	serv.Connections.Close()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	_ = httpServer.Shutdown(shutdownCtx)
	_ = redisClient.Close()
}
