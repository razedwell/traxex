package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/razedwell/traxex/services/auth/internal/domain"
	"github.com/razedwell/traxex/services/auth/internal/repository"
	"github.com/razedwell/traxex/shared/traxerr"
	"golang.org/x/crypto/bcrypt"
)

// SessionStore is a port over Redis so the service stays testable.
type SessionStore interface {
	Save(ctx context.Context, s domain.Session, ttl time.Duration) error
	Get(ctx context.Context, refresh string) (domain.Session, error)
	Delete(ctx context.Context, refresh string) error
}

type AuthService struct {
	repo       repository.UserRepository
	sessions   SessionStore
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	bcryptCost int
}

func NewAuthService(r repository.UserRepository, s SessionStore, secret []byte) *AuthService {
	return &AuthService{
		repo:       r,
		sessions:   s,
		jwtSecret:  secret,
		accessTTL:  15 * time.Minute,
		refreshTTL: 7 * 24 * time.Hour,
		bcryptCost: 12,
	}
}

func (a *AuthService) Register(ctx context.Context, email, password string) (domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), a.bcryptCost)
	if err != nil {
		return domain.User{}, traxerr.Wrap(traxerr.CodeInternal, "hash password", err)
	}
	return a.repo.Create(ctx, domain.User{
		Email:        email,
		PasswordHash: string(hash),
	})
}

func (a *AuthService) Login(ctx context.Context, email, password string) (access, refresh string, expiresIn int64, err error) {
	u, err := a.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", 0, domain.ErrInvalidCredentials()
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", "", 0, domain.ErrInvalidCredentials()
	}
	return a.issueTokens(ctx, u.ID)
}

func (a *AuthService) ValidateToken(ctx context.Context, access string) (string, error) {
	tok, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrTokenInvalid()
		}
		return a.jwtSecret, nil
	})
	if err != nil || !tok.Valid {
		return "", domain.ErrTokenInvalid()
	}

	claims := tok.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)
	if sub == "" {
		return "", domain.ErrTokenInvalid()
	}
	return sub, nil
}

func (a *AuthService) Refresh(ctx context.Context, refresh string) (string, string, int64, error) {
	sesh, err := a.sessions.Get(ctx, refresh)
	if err != nil {
		return "", "", 0, domain.ErrTokenInvalid()
	}

	err = a.sessions.Delete(ctx, refresh)
	if err != nil {
		return "", "", 0, domain.ErrTokenInvalid()
	}

	return a.issueTokens(ctx, sesh.UserID)
}

func (a *AuthService) Logout(ctx context.Context, refresh string) error {
	return a.sessions.Delete(ctx, refresh)
}

func (a *AuthService) issueTokens(ctx context.Context, userID string) (access, refresh string, expiresIn int64, err error) {
	now := time.Now()
	claims := jwt.MapClaims(map[string]interface{}{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(a.accessTTL).Unix(),
	})

	access, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(a.jwtSecret)
	if err != nil {
		return "", "", 0, traxerr.Wrap(traxerr.CodeInternal, "sign token", err)
	}

	refresh = uuid.NewString()

	sesh := domain.Session{
		UserID:       userID,
		RefreshToken: refresh,
		CreatedAt:    now,
		ExpiresAt:    now.Add(a.refreshTTL),
	}
	if err := a.sessions.Save(ctx, sesh, a.refreshTTL); err != nil {
		return "", "", 0, traxerr.Wrap(traxerr.CodeInternal, "save session", err)
	}

	return access, refresh, int64(a.accessTTL.Seconds()), nil
}
