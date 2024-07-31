package user

import (
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/cache"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/client/db"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/repository"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	logsRepository repository.LogRepository
	cache          cache.UserCache
	txManager      db.TxManager
}

// NewService creates a new user service.
func NewService(
	userRepository repository.UserRepository,
	logsRepository repository.LogRepository,
	cache cache.UserCache,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		userRepository: userRepository,
		logsRepository: logsRepository,
		cache:          cache,
		txManager:      txManager,
	}
}
