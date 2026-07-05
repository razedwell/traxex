package handler

import (
	"context"

	"github.com/razedwell/traxex/services/auth/internal/domain"
)

type AuthUseCase interface {
	Register(ctx context.Context, email, password string) (domain.User, error)
	Login(ctx context.Context, email, password string) (string, string, int64, error)
	Logout(ctx context.Context, refresh string) error
	ValidateToken(ctx context.Context, access string) (string, error)
	RefreshToken(ctx context.Context, refresh string) (string, string, int64, error)
}
