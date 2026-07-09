package service

import (
	"context"
	"testing"

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
