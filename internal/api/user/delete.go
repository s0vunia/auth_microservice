package user

import (
	"context"

	desc "github.com/s0vunia/auth_microservice/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete deletes a user.
func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
