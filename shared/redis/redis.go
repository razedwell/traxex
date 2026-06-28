package redis

import (
	"context"
	"fmt"

	"github.com/razedwell/traxex/shared/config"
	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	opt := &redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	}
	rdb := redis.NewClient(opt)

	if rdb == nil {
		return nil, fmt.Errorf("Error: Redis client is nil")
	}

	res := rdb.Ping(ctx)
	if res.Err() != nil {
		return nil, fmt.Errorf("Redis client ping error: %w-%v", res.Err(), res.Val())
	}

	return rdb, nil
}

func HealthCheck(ctx context.Context, rdb *redis.Client) error {
	if res := rdb.Ping(ctx); res.Err() != nil {
		res.Result()
		return fmt.Errorf("Redis health check failed: %w-%v", res.Err(), res.Val())
	}
	return nil
}
