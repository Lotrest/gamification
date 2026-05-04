package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	"cdek/platform/gamification/internal/application"
	postgresrepo "cdek/platform/gamification/internal/infrastructure/postgres"
	grpcserver "cdek/platform/gamification/internal/presentation/grpc"
	gamificationv1 "cdek/platform/shared/contracts/gamification/v1"
	"cdek/platform/shared/contracts/transport/jsoncodec"
)

func main() {
	jsoncodec.Register()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	address := env("GAMIFICATION_SERVICE_ADDRESS", ":50052")

	pool, err := openPostgresPool()
	if err != nil {
		logger.Error("failed to connect postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("failed to listen", "error", err, "address", address)
		os.Exit(1)
	}

	service := application.NewService(postgresrepo.NewRepository(pool))
	server := grpc.NewServer()
	gamificationv1.RegisterGamificationServiceServer(server, grpcserver.New(service))

	logger.Info("gamification-service started", "address", address)

	if err := server.Serve(listener); err != nil {
		logger.Error("grpc serve failed", "error", err)
		os.Exit(1)
	}
}

func openPostgresPool() (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env("POSTGRES_USER", "postgres"),
		env("POSTGRES_PASSWORD", "1234"),
		env("POSTGRES_HOST", "127.0.0.1"),
		env("POSTGRES_PORT", "5432"),
		env("POSTGRES_DB", "postgres"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
