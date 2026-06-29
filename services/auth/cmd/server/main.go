package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/razedwell/traxex/shared/config"
	"github.com/razedwell/traxex/shared/obs"
	"github.com/razedwell/traxex/shared/postgres"
	"github.com/razedwell/traxex/shared/redis"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Config loaded successfully")

	logger, err := obs.InitLogger(cfg.Logging)
	if err != nil {
		log.Fatalf("Failed to initialized logger: %v", err)
	}
	logger.Info("Zap logger initialized successfully")

	tracerShutdown, err := obs.InitTracer(ctx, "auth-service")
	if err != nil {
		log.Fatalf("Failed to initialize tracer for auth-service: %v", err)
	}
	defer func(func(context.Context) error) {
		logger.Info("Starting tracer shutdown...")
		tracerShutdown(context.Background())
	}(tracerShutdown)

	tracer := otel.Tracer("auth-service")
	ctx, span := tracer.Start(ctx, "smoke-test")
	span.End()

	pgpool, err := postgres.NewPool(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to create postgres pool: %v", err)
	}
	defer pgpool.Close()

	if err := postgres.HealthCheck(ctx, pgpool); err != nil {
		log.Fatalf("Failed to health check postgres: %v", err)
	}
	logger.Info("PostgreSQL pool created and healthy")

	rdb, err := redis.NewClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	if err := redis.HealthCheck(ctx, rdb); err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	logger.Info("Redis client initialized successfully")

	select {
	case <-ctx.Done():
		logger.Info("Context done chan closed", obs.TraceFields(ctx)...)
		stop()
		return
	case <-make(chan int):
		logger.Debug("Blocking channel is sending signals...")
	}
}
