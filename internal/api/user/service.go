package user

import (
	"github.com/s0vunia/auth_microservice/internal/service"
	desc "github.com/s0vunia/auth_microservice/pkg/user_v1"
)

// Implementation represents a user API implementation.
type Implementation struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

// NewImplementation creates a new user API implementation.
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
