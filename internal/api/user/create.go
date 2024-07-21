package user

import (
	"context"
	"log"

	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/converter"
	desc "github.com/s0vunia/auth_microservices_course_boilerplate/pkg/auth_v1"
)

// Create creates a new user.
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.ToUserCreateFromDesc(req.GetUserCreate()))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted note with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
