package user

import (
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/service"
	desc "github.com/s0vunia/auth_microservices_course_boilerplate/pkg/auth_v1"
)

// Implementation represents a user API implementation.
type Implementation struct {
	desc.UnimplementedAuthV1Server
	userService service.UserService
}

// NewImplementation creates a new user API implementation.
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
