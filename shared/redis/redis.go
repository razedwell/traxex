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

	res := rdb.Ping(ctx)
	if res.Err() != nil {
		return nil, fmt.Errorf("redis client ping error: %w", res.Err())
	}

	return rdb, nil
}

func HealthCheck(ctx context.Context, rdb *redis.Client) error {
	if res := rdb.Ping(ctx); res.Err() != nil {
		return fmt.Errorf("redis health check failed: %w", res.Err())
	}
	return nil
}
