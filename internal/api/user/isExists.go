package user

import (
	"context"

	desc "github.com/s0vunia/auth_microservices_course_boilerplate/pkg/auth_v1"
)

// IsExists checks if a user exists.
func (i *Implementation) IsExists(ctx context.Context, req *desc.IsExistsRequest) (*desc.IsExistsResponse, error) {
	exists, err := i.userService.IsExists(ctx, req.GetIds())
	if err != nil {
		return nil, err
	}

	return &desc.IsExistsResponse{
		Exists: exists,
	}, nil
}
