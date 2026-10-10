package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lama189/soundwave-platform/catalog-service/internal/config"
	deliveryGrpc "github.com/Lama189/soundwave-platform/catalog-service/internal/delivery/grpc"
	"github.com/Lama189/soundwave-platform/catalog-service/internal/repository/postgres"
	"github.com/Lama189/soundwave-platform/catalog-service/internal/service/artist"
	catalogpb "github.com/Lama189/soundwave-platform/gen/go/catalog/v1"
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

	log.Println("Database connection established and migrations applied successfully")

	artistRepo := postgres.NewArtistRepository(pgPool)
	artistService := artist.NewService(artistRepo)

	grpcServer := grpc.NewServer()
	artistGrpcHandler := deliveryGrpc.NewArtistGRPCServer(artistService)
	catalogpb.RegisterArtistServiceServer(grpcServer, artistGrpcHandler)

	lis, err := net.Listen("tcp", cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", cfg.GRPC.Port, err)
	}

	go func() {
		log.Printf("Starting catalog-service gRPC server on %s", cfg.GRPC.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down catalog-service gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Server exited properly")
}
