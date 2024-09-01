package user

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/converter"
	"github.com/s0vunia/auth_microservice/internal/logger"
	desc "github.com/s0vunia/auth_microservice/pkg/user_v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Update updates a user.
func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	logger.Info("Updating user...", zap.Int64("id", req.GetId()))
	err := i.userService.Update(ctx, req.GetId(), converter.ToUserUpdateFromDesc(req.GetUserUpdate()))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
