package access

import (
	"github.com/s0vunia/auth_microservice/internal/cache"
	"github.com/s0vunia/auth_microservice/internal/repository"
	"github.com/s0vunia/auth_microservice/internal/service"
	"github.com/s0vunia/platform_common/pkg/db"
)

type serv struct {
	logsRepository       repository.LogRepository
	cache                cache.UserCache
	txManager            db.TxManager
	authPrefix           string
	accessTokenSecretKey string
}

// NewService creates a new access service.
func NewService(
	logsRepository repository.LogRepository,
	cache cache.UserCache,
	txManager db.TxManager,
	authPrefix string,
	accessTokenSecretKey string,
) service.AccessService {
	return &serv{
		logsRepository:       logsRepository,
		cache:                cache,
		txManager:            txManager,
		authPrefix:           authPrefix,
		accessTokenSecretKey: accessTokenSecretKey,
	}
}
