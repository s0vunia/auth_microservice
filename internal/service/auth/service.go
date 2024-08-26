package auth

import (
	"time"

	"github.com/s0vunia/auth_microservice/internal/cache"
	"github.com/s0vunia/auth_microservice/internal/repository"
	"github.com/s0vunia/auth_microservice/internal/service"
	"github.com/s0vunia/platform_common/pkg/db"
)

type serv struct {
	userRepository         repository.UserRepository
	cache                  cache.UserCache
	txManager              db.TxManager
	refreshTokenSecretKey  string
	refreshTokenExpiration time.Duration
	accessTokenSecretKey   string
	accessTokenExpiration  time.Duration
}

// NewService creates a new auth service.
func NewService(
	userRepository repository.UserRepository,
	cache cache.UserCache,
	txManager db.TxManager,
	refreshTokenSecretKey string,
	refreshTokenExpiration time.Duration,
	accessTokenSecretKey string,
	accessTokenExpiration time.Duration,
) service.AuthService {
	return &serv{
		userRepository:         userRepository,
		cache:                  cache,
		txManager:              txManager,
		refreshTokenSecretKey:  refreshTokenSecretKey,
		refreshTokenExpiration: refreshTokenExpiration,
		accessTokenSecretKey:   accessTokenSecretKey,
		accessTokenExpiration:  accessTokenExpiration,
	}
}
