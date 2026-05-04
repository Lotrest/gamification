package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cdek/platform/bff/internal/app"
	"cdek/platform/bff/internal/middleware"
	gamificationv1 "cdek/platform/shared/contracts/gamification/v1"
	"cdek/platform/shared/contracts/transport/jsoncodec"
	userv1 "cdek/platform/shared/contracts/user/v1"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	jsoncodec.Register()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	tracerProvider := sdktrace.NewTracerProvider()
	defer func() {
		_ = tracerProvider.Shutdown(context.Background())
	}()

	otel.SetTracerProvider(tracerProvider)
	tracer := otel.Tracer("cdek.platform.bff")

	pool, err := openPostgresPool()
	if err != nil {
		logger.Error("failed to connect postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	userConnection, err := grpc.NewClient(
		env("USER_SERVICE_GRPC", "127.0.0.1:50051"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.ForceCodec(jsoncodec.Codec{}),
			grpc.CallContentSubtype(jsoncodec.Name),
		),
	)
	if err != nil {
		logger.Error("failed to connect to user-service", "error", err)
		os.Exit(1)
	}
	defer func() {
		_ = userConnection.Close()
	}()

	gamificationConnection, err := grpc.NewClient(
		env("GAMIFICATION_SERVICE_GRPC", "127.0.0.1:50052"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.ForceCodec(jsoncodec.Codec{}),
			grpc.CallContentSubtype(jsoncodec.Name),
		),
	)
	if err != nil {
		logger.Error("failed to connect to gamification-service", "error", err)
		os.Exit(1)
	}
	defer func() {
		_ = gamificationConnection.Close()
	}()

	registry := prometheus.NewRegistry()
	metrics := middleware.NewMetrics(registry)
	startMetricsServer(logger, env("BFF_METRICS_ADDRESS", ":9090"), registry)

	server := app.NewServer(
		logger,
		tracer,
		pool,
		userv1.NewUserServiceClient(userConnection),
		gamificationv1.NewGamificationServiceClient(gamificationConnection),
	)

	fiberApp := fiber.New()
	fiberApp.Use(recovermw.New())
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,http://127.0.0.1:5173,http://localhost:4173,http://127.0.0.1:4173,http://localhost:8080,http://127.0.0.1:8080,http://localhost:18080,http://127.0.0.1:18080",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,DELETE,OPTIONS",
	}))
	fiberApp.Use(metrics.Middleware())

	fiberApp.Get("/healthz", server.Health)

	publicAPI := fiberApp.Group("/api/v1")
	publicAPI.Post("/auth/register", server.Register)
	publicAPI.Post("/auth/login", server.Login)

	api := fiberApp.Group("/api/v1", middleware.Auth())
	api.Get("/bootstrap", server.Bootstrap)
	api.Post("/tasks/:taskId/accept", middleware.RequireRouteParam("taskId"), server.AcceptTask)
	api.Post("/tasks/:taskId/advance", middleware.RequireRouteParam("taskId"), server.AdvanceTask)
	api.Post("/rewards/:rewardId/redeem", middleware.RequireRouteParam("rewardId"), server.RedeemReward)
	api.Get("/articles", server.ListArticles)
	api.Get("/articles/:articleId", middleware.RequireRouteParam("articleId"), server.GetArticle)
	api.Post("/articles", server.CreateArticle)
	api.Post("/articles/:articleId/reactions", middleware.RequireRouteParam("articleId"), server.ToggleReaction)
	api.Post("/articles/:articleId/comments", middleware.RequireRouteParam("articleId"), server.CreateComment)
	api.Delete(
		"/articles/:articleId/comments/:commentId",
		middleware.RequireRouteParam("articleId"),
		middleware.RequireRouteParam("commentId"),
		server.DeleteComment,
	)

	address := env("BFF_ADDRESS", ":8080")
	go func() {
		logger.Info("bff started", "address", address)
		if err := fiberApp.Listen(address); err != nil {
			logger.Error("bff listen failed", "error", err)
			os.Exit(1)
		}
	}()

	waitForShutdown(logger, fiberApp)
}

func startMetricsServer(logger *slog.Logger, address string, registry *prometheus.Registry) {
	server := &http.Server{
		Addr:              address,
		Handler:           promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("metrics server started", "address", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("metrics server failed", "error", err)
		}
	}()
}

func waitForShutdown(logger *slog.Logger, fiberApp *fiber.App) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("shutdown signal received")

	shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := fiberApp.ShutdownWithContext(shutdownContext); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
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
