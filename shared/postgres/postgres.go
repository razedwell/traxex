package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/razedwell/traxex/shared/config"
)

func NewPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	pgxpoolConfig, err := buildPgxPoolConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to build pgxpool config: %w", err)
	}

	pgpool, err := pgxpool.NewWithConfig(ctx, pgxpoolConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create pgxpool: %w", err)
	}

	if err := pgpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	return pgpool, nil
}

func HealthCheck(ctx context.Context, pgpool *pgxpool.Pool) error {
	return pgpool.Ping(ctx)
}

func buildPgxPoolConfig(cfg config.DatabaseConfig) (*pgxpool.Config, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode,
	)
	connConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse connection string: %w", err)
	}

	connConfig.ConnConfig.User = cfg.User
	connConfig.ConnConfig.Password = cfg.Password
	connConfig.ConnConfig.Database = cfg.Database
	connConfig.MaxConns = cfg.PoolMaxConns

	connConfig.ConnConfig.Host = cfg.Host
	connConfig.ConnConfig.Port = cfg.Port

	return connConfig, nil
}
