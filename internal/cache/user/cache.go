package user

import (
	"context"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/s0vunia/auth_microservice/internal/cache"
	"github.com/s0vunia/auth_microservice/internal/cache/user/converter"
	modelCache "github.com/s0vunia/auth_microservice/internal/cache/user/model"
	"github.com/s0vunia/auth_microservice/internal/model"
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

func (c *cacheImplementation) Create(ctx context.Context, user *model.User) error {
	userCache := converter.ToUserCacheFromModel(user)
	idStr := strconv.FormatInt(userCache.ID, 10)
	err := c.cl.HashSet(ctx, idStr, userCache)
	if err != nil {
		return err
	}
	return nil
}
func (c *cacheImplementation) Get(ctx context.Context, id int64) (*model.User, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := c.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, errors.New("user not found")
	}

	var user modelCache.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromCache(&user), nil
}

func (c *cacheImplementation) Delete(ctx context.Context, id int64) error {
	idStr := strconv.FormatInt(id, 10)
	err := c.cl.Del(ctx, idStr)
	if err != nil {
		return err
	}
	return nil
}
