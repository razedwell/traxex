package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/razedwell/traxex/shared/config"
)

func NewPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%d/?sslmode=%s", cfg.Host, cfg.Port, cfg.SSLMode)
	// connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%d", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode, cfg.PoolMaxConns)
	connConfig, _ := pgxpool.ParseConfig(connString)

	connConfig.ConnConfig.User = cfg.User
	connConfig.ConnConfig.Password = cfg.Password
	connConfig.ConnConfig.Database = cfg.Database
	connConfig.MaxConns = cfg.PoolMaxConns

	connConfig.ConnConfig.Host = cfg.Host
	connConfig.ConnConfig.Port = cfg.Port

	pgpool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("Failed to create pgxpool: %v", err)
	}
	return pgpool, nil
}

func HealthCheck(ctx context.Context, pgpool *pgxpool.Pool) error {
	return pgpool.Ping(ctx)
}
