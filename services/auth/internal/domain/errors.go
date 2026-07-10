package domain

import (
	"github.com/razedwell/traxex/shared/traxerr"
)

func ErrEmailTaken() error {
	return traxerr.New(traxerr.CodeConflict, "email already registered")
}

func ErrInvalidCredentials() error {
	return traxerr.New(traxerr.CodeUnauthorized, "invalid email or password")
}

func ErrUserNotFound() error {
	return traxerr.New(traxerr.CodeNotFound, "user not found")
}

func ErrTokenInvalid() error {
	return traxerr.New(traxerr.CodeUnauthorized, "token invalid or expired")
}

func ErrSessionNotFound() error {
	return traxerr.New(traxerr.CodeSessionNotFound, "session not found")
}
