package auth

import (
	"github.com/s0vunia/auth_microservice/internal/service"
	desc "github.com/s0vunia/auth_microservice/pkg/auth_v1"
)

// Implementation represents a auth API implementation.
type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewImplementation creates a new auth API implementation.
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
