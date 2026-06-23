package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/razedwell/traxex/shared/config"
)

func NewPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%d", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode, cfg.PoolMaxConns)
	pgpool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("Failed to create pgxpool: %v", err)
	}
	return pgpool, nil
}
