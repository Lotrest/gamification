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

	"cdek/platform/shared/contracts/transport/jsoncodec"
	userv1 "cdek/platform/shared/contracts/user/v1"
	"cdek/platform/user-service/internal/application"
	postgresrepo "cdek/platform/user-service/internal/infrastructure/postgres"
	grpcserver "cdek/platform/user-service/internal/presentation/grpc"
)

func main() {
	jsoncodec.Register()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	address := env("USER_SERVICE_ADDRESS", ":50051")

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
	userv1.RegisterUserServiceServer(server, grpcserver.New(service))

	logger.Info("user-service started", "address", address)

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
