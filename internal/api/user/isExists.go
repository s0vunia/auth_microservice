package user

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/logger"
	desc "github.com/s0vunia/auth_microservice/pkg/user_v1"
	"go.uber.org/zap"
)

// IsExists checks if a user exists.
func (i *Implementation) IsExists(ctx context.Context, req *desc.IsExistsRequest) (*desc.IsExistsResponse, error) {
	logger.Info("Checking if user exists...", zap.Int64s("ids", req.GetIds()))
	exists, err := i.userService.IsExists(ctx, req.GetIds())
	if err != nil {
		return nil, err
	}

	return &desc.IsExistsResponse{
		Exists: exists,
	}, nil
}
