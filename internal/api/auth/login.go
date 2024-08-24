package auth

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/logger"
	"github.com/s0vunia/auth_microservice/internal/model"
	desc "github.com/s0vunia/auth_microservice/pkg/auth_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login implements auth.Login
func (s *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	logger.Info("Logging in...", zap.Int64("id", req.Id))
	refreshToken, err := s.authService.Login(ctx, &model.UserLogin{
		ID:       req.Id,
		Password: req.Password,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &desc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
