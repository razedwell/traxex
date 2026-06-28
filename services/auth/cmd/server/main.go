package main

import (
	"context"
	"log"

	"github.com/razedwell/traxex/shared/config"
	"github.com/razedwell/traxex/shared/postgres"
	"github.com/razedwell/traxex/shared/redis"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Config loaded successfully")

	pgpool, err := postgres.NewPool(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to create postgres pool: %v", err)
	}
	defer pgpool.Close()

	if err := postgres.HealthCheck(ctx, pgpool); err != nil {
		log.Fatalf("Failed to health check postgres: %v", err)
	}

	log.Println("PostgreSQL pool created and healthy")

	rdb, err := redis.NewClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	if err := redis.HealthCheck(ctx, rdb); err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	log.Println("Redis client initialized successfully")

	select {
	case <-ctx.Done():
		log.Println("Context done")
		return
	default:
		log.Println("Context still running")
	}
}
