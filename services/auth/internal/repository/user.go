package repository

import (
	"context"

	"github.com/razedwell/traxex/services/auth/internal/domain"
)

// UserRepository is the ONLY thing the service layer knows about persistence.
type UserRepository interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByID(ctx context.Context, id string) (domain.User, error)
}
