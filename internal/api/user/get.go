package user

import (
	"context"
	"log"

	"github.com/s0vunia/auth_microservice/internal/converter"
	desc "github.com/s0vunia/auth_microservice/pkg/auth_v1"
)

// Get gets a user.
func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	userObj, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, name: %s, email: %s, role: %d, created_at: %v, updated_at: %v\n", userObj.ID, userObj.Info.Name, userObj.Info.Email, userObj.Info.Role, userObj.CreatedAt, userObj.UpdatedAt)

	return &desc.GetResponse{
		User: converter.ToUserFromService(userObj),
	}, nil
}
