package handler

// import (
// 	"context"
// 	"errors"
// 	"net/http"

// 	authv1 "github.com/razedwell/traxex/proto/gen/go/auth/v1"
// 	"github.com/razedwell/traxex/shared/traxerr"
// 	"google.golang.org/protobuf/types/known/timestamppb"
// )

// type HTTPHandler struct {
// 	authv1.UnimplementedAuthServiceServer
// 	svc AuthUseCase
// }

// func (h *HTTPHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
// 	u, err := h.svc.Register(ctx, req.GetEmail(), req.GetPassword())
// 	if err != nil {
// 		return nil, toHTTPError(err)
// 	}
// 	return &authv1.RegisterResponse{
// 		User: &authv1.User{
// 			Id:        u.ID,
// 			Email:     u.Email,
// 			CreatedAt: timestamppb.New(u.CreatedAt),
// 		}}, nil
// }

// func (h *HTTPHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
// 	access, refresh, exp, err := h.svc.Login(ctx, req.GetEmail(), req.GetPassword())
// 	if err != nil {
// 		return nil, toHTTPError(err)
// 	}
// 	return &authv1.LoginResponse{
// 		AccessToken:  access,
// 		RefreshToken: refresh,
// 		ExpiresIn:    exp,
// 	}, nil
// }

// func (h *HTTPHandler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
// 	err := h.svc.Logout(ctx, req.GetRefreshToken())
// 	if err != nil {
// 		return &authv1.LogoutResponse{Success: false}, toHTTPError(err)
// 	}
// 	return &authv1.LogoutResponse{Success: true}, nil
// }

// func (h *HTTPHandler) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
// 	userID, err := h.svc.ValidateToken(ctx, req.GetAccessToken())
// 	if err != nil {
// 		return nil, toHTTPError(err)
// 	}
// 	return &authv1.ValidateTokenResponse{
// 		UserId: userID,
// 	}, nil
// }

// func (h *HTTPHandler) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
// 	access, refresh, exp, err := h.svc.RefreshToken(ctx, req.GetRefreshToken())
// 	if err != nil {
// 		return nil, toHTTPError(err)
// 	}
// 	return &authv1.RefreshTokenResponse{
// 		AccessToken:  access,
// 		RefreshToken: refresh,
// 		ExpiresIn:    exp,
// 	}, nil
// }

// func toHTTPError(e error) error {
// 	http.Error()
// 	var te *traxerr.Traxerr
// 	if errors.As(e, te) {
// 		return traxerr.HTTPStatus(te.Code)
// 	}
// }
