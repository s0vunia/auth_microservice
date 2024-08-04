package user

import (
	"github.com/s0vunia/auth_microservice/internal/cache"
	cacheCl "github.com/s0vunia/platform_common/pkg/cache"
)

type cacheImplementation struct {
	cl cacheCl.Client
}

func NewCache(cl cacheCl.Client) cache.UserCache {
	return &cacheImplementation{
		cl: cl,
	}
}
