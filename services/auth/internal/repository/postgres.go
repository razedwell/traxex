package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/razedwell/traxex/services/auth/internal/domain"
)

type PostgresUserRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepo(pool *pgxpool.Pool) *PostgresUserRepo {
	return &PostgresUserRepo{pool}
}

func (r *PostgresUserRepo) Create(ctx context.Context, u domain.User) (domain.User, error) {
	const q = `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	err := r.pool.QueryRow(ctx, q, u.Email, u.PasswordHash).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, domain.ErrEmailTaken()
		}
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	const q = `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	var u domain.User
	err := r.pool.QueryRow(ctx, q, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound()
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by email: %w", err)
	}
	return u, nil
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	const q = `SELECT id, email, password_hash, created_at FROM users WHERE id = $1`
	var u domain.User
	err := r.pool.QueryRow(ctx, q, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound()
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by id: %w", err)
	}
	return u, nil
}
