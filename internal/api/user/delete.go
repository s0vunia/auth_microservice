package user

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/logger"
	desc "github.com/s0vunia/auth_microservice/pkg/user_v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete deletes a user.
func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	logger.Info("Deleting user...", zap.Int64("id", req.GetId()))
	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
