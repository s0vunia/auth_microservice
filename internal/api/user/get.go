package user

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/converter"
	"github.com/s0vunia/auth_microservice/internal/logger"
	desc "github.com/s0vunia/auth_microservice/pkg/user_v1"
	"go.uber.org/zap"
)

// Get gets a user.
func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	logger.Info("Getting user...", zap.Int64("id", req.GetId()))
	userObj, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		User: converter.ToUserFromService(userObj),
	}, nil
}
