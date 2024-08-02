package user

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/converter"
	desc "github.com/s0vunia/auth_microservice/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Update updates a user.
func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, req.GetId(), converter.ToUserUpdateFromDesc(req.GetUserUpdate()))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
