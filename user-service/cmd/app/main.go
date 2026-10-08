package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	userpb "github.com/Lama189/soundwave-platform/gen/go/user/v1"
	"github.com/Lama189/soundwave-platform/user-service/internal/config"
	deliveryGrpc "github.com/Lama189/soundwave-platform/user-service/internal/delivery/grpc/v1"
	"github.com/Lama189/soundwave-platform/user-service/internal/infrastructure/hasher"
	"github.com/Lama189/soundwave-platform/user-service/internal/infrastructure/jwt"
	"github.com/Lama189/soundwave-platform/user-service/internal/repository/postgres"
	"github.com/Lama189/soundwave-platform/user-service/internal/repository/redis"
	"github.com/Lama189/soundwave-platform/user-service/internal/service/user"
	"google.golang.org/grpc"
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

	grpcServer := grpc.NewServer()
	userGrpcHandler := deliveryGrpc.NewUserGRPCServer(userService)
	userpb.RegisterUserServiceServer(grpcServer, userGrpcHandler)

	lis, err := net.Listen("tcp", cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", cfg.GRPC.Port, err)
	}

	go func() {
		log.Printf("Starting gRPC server on %s", cfg.GRPC.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Server exited properly")
}
