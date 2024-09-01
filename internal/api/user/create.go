package user

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/converter"
	"github.com/s0vunia/auth_microservice/internal/logger"
	desc "github.com/s0vunia/auth_microservice/pkg/user_v1"
	"go.uber.org/zap"
)

// Create creates a new user.
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	logger.Info(
		"Creating user...",
		zap.String("name", req.GetUserCreate().GetInfo().GetName()),
	)
	if req.GetUserCreate().GetPassword() != req.GetUserCreate().GetPasswordConfirm() {
		return nil, PasswordConfirmErr
	}

	id, err := i.userService.Create(ctx, converter.ToUserCreateFromDesc(req.GetUserCreate()))
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
