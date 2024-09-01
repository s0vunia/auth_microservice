package user

import (
	"github.com/s0vunia/auth_microservice/internal/cache"
	"github.com/s0vunia/auth_microservice/internal/repository"
	"github.com/s0vunia/auth_microservice/internal/service"
	"github.com/s0vunia/platform_common/pkg/db"
)

type serv struct {
	userRepository repository.UserRepository
	cache          cache.UserCache
	txManager      db.TxManager
}

// NewService creates a new user service.
func NewService(
	userRepository repository.UserRepository,
	cache cache.UserCache,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		userRepository: userRepository,
		cache:          cache,
		txManager:      txManager,
	}
}
