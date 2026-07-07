package handler

import (
	"context"
	"errors"

	authv1 "github.com/razedwell/traxex/proto/gen/go/auth/v1"
	"github.com/razedwell/traxex/services/auth/internal/service"
	"github.com/razedwell/traxex/shared/traxerr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	authv1.UnimplementedAuthServiceServer
	svc AuthUseCase
}

// compile-time checks
var _ AuthUseCase = (*service.AuthService)(nil)
var _ authv1.AuthServiceServer = (*GRPCHandler)(nil)

func NewGRPCHandler(svc AuthUseCase) *GRPCHandler {
	return &GRPCHandler{svc: svc}
}

func (h *GRPCHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	u, err := h.svc.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	return &authv1.RegisterResponse{
		User: &authv1.User{
			Id:        u.ID,
			Email:     u.Email,
			CreatedAt: timestamppb.New(u.CreatedAt),
		}}, nil
}

func (h *GRPCHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	access, refresh, exp, err := h.svc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	return &authv1.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    exp,
	}, nil
}

func (h *GRPCHandler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	err := h.svc.Logout(ctx, req.GetRefreshToken())
	if err != nil {
		return &authv1.LogoutResponse{Success: false}, toGRPCStatus(err)
	}
	return &authv1.LogoutResponse{Success: true}, nil
}

func (h *GRPCHandler) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	userID, err := h.svc.ValidateToken(ctx, req.GetAccessToken())
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	return &authv1.ValidateTokenResponse{
		UserId: userID,
	}, nil
}

func (h *GRPCHandler) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	access, refresh, exp, err := h.svc.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	return &authv1.RefreshTokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    exp,
	}, nil
}

// transforms an error into a gRPC Status (error/failure).
func toGRPCStatus(e error) error {
	var te *traxerr.Traxerr
	if errors.As(e, &te) {
		return status.Error(traxerr.GRPCCode(te.Code), te.Message)
	}
	return status.Error(codes.Internal, "internal server error")
}
