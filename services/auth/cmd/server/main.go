package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/razedwell/traxex/shared/config"
	"github.com/razedwell/traxex/shared/kafka"
	"github.com/razedwell/traxex/shared/obs"
	"github.com/razedwell/traxex/shared/postgres"
	"github.com/razedwell/traxex/shared/redis"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
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
	defer logger.Sync()
	logger.Info("Zap logger initialized successfully")

	tracerShutdown, err := obs.InitTracer(ctx, "auth-service")
	if err != nil {
		logger.Fatal("Failed to initialize tracer for auth-service: %v", zap.Error(err))
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tracerShutdown(shutdownCtx); err != nil {
			logger.Error("Tracer shutdown failed", zap.Error(err))
		}
		logger.Info("Starting tracer shutdown...")
	}()

	tracer := otel.Tracer("auth-service")
	ctx, span := tracer.Start(ctx, "smoke-test")

	pgpool, err := postgres.NewPool(ctx, cfg.Database)
	if err != nil {
		logger.Fatal("Failed to create postgres pool: %v", zap.Error(err))
	}
	defer pgpool.Close()

	if err := postgres.HealthCheck(ctx, pgpool); err != nil {
		logger.Fatal("Failed to health check postgres: %v", zap.Error(err))
	}
	logger.Info("PostgreSQL pool created and healthy")

	rdb, err := redis.NewClient(ctx, cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to connect to Redis: %v", zap.Error(err))
	}

	if err := redis.HealthCheck(ctx, rdb); err != nil {
		logger.Fatal("Failed to ping Redis: %v", zap.Error(err))
	}

	logger.Info("Redis client initialized successfully")

	topic := "phase1.smoke"

	producer := kafka.NewProducer(cfg.Kafka, topic)
	defer producer.Close()
	logger.Info("Kafka producer initialized successfully")

	consumer := kafka.NewConsumer(cfg.Kafka, topic)
	defer consumer.Close()
	logger.Info("Kafka consumer initialized successfully")

	err = producer.Publish(ctx, []byte("Hi"), []byte("Kafka Apache"))
	if err != nil {
		logger.Fatal("Kafka producer publish failed", obs.TraceFields(ctx, zap.Error(err))...)
	}

	msg, err := consumer.Read(ctx)
	if err != nil {
		logger.Fatal("Kafka consumer read failed", obs.TraceFields(ctx, zap.Error(err))...)
	}
	logger.Info("From Kafka got message", zap.String("value", string(msg.Value)))

	span.End()

	select {
	case <-ctx.Done():
		logger.Info("Context done chan closed", obs.TraceFields(ctx)...)
		stop()
		return
	case <-make(chan int):
		logger.Debug("Blocking channel is sending signals...")
	}
}
