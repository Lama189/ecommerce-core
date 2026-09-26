package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lama189/ecommerce-core/user-service/internal/config"
	deliveryHttp "github.com/Lama189/ecommerce-core/user-service/internal/delivery/http"
	"github.com/Lama189/ecommerce-core/user-service/internal/infrastructure/hasher"
	"github.com/Lama189/ecommerce-core/user-service/internal/infrastructure/jwt"
	"github.com/Lama189/ecommerce-core/user-service/internal/repository/postgres"
	"github.com/Lama189/ecommerce-core/user-service/internal/repository/redis"
	"github.com/Lama189/ecommerce-core/user-service/internal/service/user"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()

	pgPool, err := postgres.NewPool(ctx, postgres.Config{
		DSN: cfg.Postgres.URL,
	})
	if err != nil {
		log.Fatalf("failed to init postgres: %v", err)
	}
	defer pgPool.Close()

	if err := postgres.RunMigrations(ctx, pgPool); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	redisClient, err := redis.NewClient(ctx, redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		log.Fatalf("failed to init redis: %v", err)
	}
	defer redisClient.Close()

	log.Println("Database and cache connections established successfully")

	pwdHasher := hasher.NewBcryptHahser(0)
	jwtManager := jwt.NewJWTManager(cfg.JWT.SecretKey, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	userRepo := postgres.NewUserRepository(pgPool)
	userCache := redis.NewCacheRepository(redisClient, cfg.Redis.UserTTL, cfg.JWT.RefreshTTL)

	userService := user.NewService(userRepo, userCache, pwdHasher, jwtManager)

	userHandler := deliveryHttp.NewUserHandler(userService)
	router := deliveryHttp.NewRouter(userHandler, jwtManager)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
