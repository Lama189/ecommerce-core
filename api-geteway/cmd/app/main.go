package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lama189/ecommerce-core/api-geteway/internal/client"
	userClient "github.com/Lama189/ecommerce-core/api-geteway/internal/client/user"
	"github.com/Lama189/ecommerce-core/api-geteway/internal/config"
	deliveryHttp "github.com/Lama189/ecommerce-core/api-geteway/internal/delivery/http"
	"github.com/Lama189/ecommerce-core/api-geteway/internal/infrastructure/jwt"
	"github.com/Lama189/ecommerce-core/api-geteway/internal/infrastructure/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)
	log.Info(
		"starting API Gateway",
		slog.String("env", cfg.Env),
		slog.String("port", cfg.HTTP.Port),
	)

	jwtValidator := jwt.NewValidator(cfg.JWT.SecretKey)

	clientMgr := client.NewManager(log)

	initCtx, cancelInit := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelInit()

	userConn, err := clientMgr.GetConn(initCtx, cfg.Services.User.Addr)
	if err != nil {
		log.Error(
			"failed to connect to user-service",
			slog.String("addr", cfg.Services.User.Addr),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	userSvcClient := userClient.NewClient(userConn, cfg.Services.User.Timeout)
	userHandler := deliveryHttp.NewUserHandler(userSvcClient)
	router := deliveryHttp.NewRouter(userHandler, jwtValidator)

	srv := &http.Server{
		Addr:         cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	go func() {
		log.Info("HTTP server listening", slog.String("addr", cfg.HTTP.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down API Gateway...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP server forced to shutdown", slog.String("error", err.Error()))
	}

	if err := clientMgr.Close(); err != nil {
		log.Error("failed to close gRPC connections", slog.String("error", err.Error()))
	}

	log.Info("API Gateway stopped gracefully")
}
