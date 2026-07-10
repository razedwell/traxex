package service

import (
	"context"
	"testing"
	"time"

	"github.com/razedwell/traxex/services/auth/internal/domain"
	"github.com/razedwell/traxex/services/auth/internal/mocks"

	"github.com/razedwell/traxex/shared/traxerr"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), 4)

	tests := []struct {
		name    string
		pswd    string
		setup   func(r *mocks.MockUserRepository, s *mocks.MockSessionStore)
		wantErr traxerr.Code
	}{
		{
			name: "success",
			pswd: "correct-password",
			setup: func(r *mocks.MockUserRepository, s *mocks.MockSessionStore) {
				r.EXPECT().GetByEmail(mock.Anything, "a@b.cd").
					Return(domain.User{ID: "u1", PasswordHash: string(hash)}, nil)
				s.EXPECT().Save(mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name: "fail",
			pswd: "wrong pswd",
			setup: func(r *mocks.MockUserRepository, s *mocks.MockSessionStore) {
				r.EXPECT().GetByEmail(mock.Anything, "a@b.cd").
					Return(domain.User{ID: "u2", PasswordHash: string(hash)}, nil)
			},
			wantErr: traxerr.CodeUnauthorized,
		},
		{
			name: "unknown email returns the SAME unauthorized (enumeration defense)",
			pswd: "whatever",
			setup: func(r *mocks.MockUserRepository, s *mocks.MockSessionStore) {
				r.EXPECT().GetByEmail(mock.Anything, "a@b.cd").
					Return(domain.User{}, domain.ErrUserNotFound())
			},
			wantErr: traxerr.CodeUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, sesh := mocks.NewMockUserRepository(t), mocks.NewMockSessionStore(t)
			tt.setup(repo, sesh)
			svc := NewAuthService(repo, sesh, []byte("test-secret"))

			_, _, _, err := svc.Login(context.Background(), "a@b.cd", tt.pswd)

			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.True(t, traxerr.IsCode(err, tt.wantErr), "got: %v", err)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), 4)

	tests := []struct {
		name    string
		pass    string
		setup   func(r *mocks.MockUserRepository, s *mocks.MockSessionStore)
		wantErr traxerr.Code
	}{
		{
			name: "success",
			pass: "correct-password",
			setup: func(r *mocks.MockUserRepository, s *mocks.MockSessionStore) {
				r.EXPECT().Create(mock.Anything, mock.Anything).
					Return(domain.User{ID: "u1", Email: "a@b.cd", PasswordHash: string(hash)}, nil)
			},
		},
		{
			name: "fail",
			pass: "email taken",
			setup: func(r *mocks.MockUserRepository, s *mocks.MockSessionStore) {
				r.EXPECT().Create(mock.Anything, mock.Anything).
					Return(domain.User{}, domain.ErrEmailTaken())
			},
			wantErr: traxerr.CodeConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, sesh := mocks.NewMockUserRepository(t), mocks.NewMockSessionStore(t)
			tt.setup(repo, sesh)
			svc := NewAuthService(repo, sesh, []byte("test-secret"))

			_, err := svc.Register(context.Background(), "a@b.cd", tt.pass)

			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.True(t, traxerr.IsCode(err, tt.wantErr), "got: %v", err)
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		setup   func(r *mocks.MockUserRepository, s *mocks.MockSessionStore)
		wantErr traxerr.Code
	}{
		{
			name:  "valid",
			input: "valid token",
			setup: func(r *mocks.MockUserRepository, s *mocks.MockSessionStore) {
				s.EXPECT().Get(mock.Anything, mock.Anything).
					Return(domain.Session{UserID: "u1", RefreshToken: "valid token", CreatedAt: time.Now(), ExpiresAt: time.Now().Add(15 * time.Minute)}, nil)
				s.EXPECT().Delete(mock.Anything, mock.Anything).Return(nil)
				s.EXPECT().Save(mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:  "invalid",
			input: "invalid token",
			setup: func(r *mocks.MockUserRepository, s *mocks.MockSessionStore) {
				s.EXPECT().Get(mock.Anything, mock.Anything).
					Return(domain.Session{}, domain.ErrTokenInvalid())
			},
			wantErr: traxerr.CodeUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, sesh := mocks.NewMockUserRepository(t), mocks.NewMockSessionStore(t)
			tt.setup(repo, sesh)
			svc := NewAuthService(repo, sesh, []byte("test-secret"))

			_, _, _, err := svc.RefreshToken(context.Background(), mock.Anything)
			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.True(t, traxerr.IsCode(err, traxerr.CodeUnauthorized), "got %v: ", err)
			}
		})
	}
}
