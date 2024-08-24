package auth

import (
	"context"
	"errors"

	"github.com/s0vunia/auth_microservice/internal/logger"
	"github.com/s0vunia/auth_microservice/internal/service/auth"
	desc "github.com/s0vunia/auth_microservice/pkg/auth_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRefreshToken implements auth.GetRefreshToken
func (s *Implementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	logger.Info("Getting refresh token...", zap.String("refresh_token", req.RefreshToken))
	refreshToken, err := s.authService.GetRefreshToken(ctx, req.RefreshToken)
	if errors.Is(err, auth.ErrInvalidRefreshToken) {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &desc.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
