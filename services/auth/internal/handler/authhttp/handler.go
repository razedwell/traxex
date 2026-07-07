package authhttp

import (
	"context"

	"github.com/razedwell/traxex/services/auth/internal/handler"
	"github.com/razedwell/traxex/shared/traxerr"
)

type HTTPHandler struct {
	svc handler.AuthUseCase
}

func NewHandler(svc handler.AuthUseCase) *HTTPHandler {
	return &HTTPHandler{svc: svc}
}

func (h *HTTPHandler) Register(ctx context.Context, req RegisterRequestObject) (RegisterResponseObject, error) {
	u, err := h.svc.Register(ctx, string(req.Body.Email), req.Body.Password)
	if err != nil {
		if traxerr.IsCode(err, traxerr.CodeConflict) {
			return Register409Response{}, nil
		}
		return nil, err
	}
	return Register201JSONResponse{
		Id:        &u.ID,
		Email:     &u.Email,
		CreatedAt: &u.CreatedAt,
	}, nil
}

func (h *HTTPHandler) Login(ctx context.Context, req LoginRequestObject) (LoginResponseObject, error) {
	access, refresh, exp, err := h.svc.Login(ctx, string(req.Body.Email), req.Body.Password)
	if err != nil {
		if traxerr.IsCode(err, traxerr.CodeInvalidInput) {
			return Login401Response{}, nil
		}
		return nil, err
	}
	return Login200JSONResponse{
		AccessToken:  &access,
		RefreshToken: &refresh,
		ExpiresIn:    &exp,
	}, nil
}

func (h *HTTPHandler) Logout(ctx context.Context, req LogoutRequestObject) (LogoutResponseObject, error) {
	var success bool
	if err := h.svc.Logout(ctx, req.Body.RefreshToken); err != nil {
		if traxerr.IsCode(err, traxerr.CodeSessionNotFound) {
			return Logout204Response{}, nil
		}
		return nil, err
	}

	success = true
	return Logout200JSONResponse{
		Success: &success,
	}, nil
}
