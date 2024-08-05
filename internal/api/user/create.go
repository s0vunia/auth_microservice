package user

import (
	"context"
	"log"

	"github.com/s0vunia/auth_microservice/internal/converter"
	desc "github.com/s0vunia/auth_microservice/pkg/auth_v1"
)

// Create creates a new user.
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if req.GetUserCreate().GetPassword() != req.GetUserCreate().GetPasswordConfirm() {
		return nil, PasswordConfirmErr
	}

	id, err := i.userService.Create(ctx, converter.ToUserCreateFromDesc(req.GetUserCreate()))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted user with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
