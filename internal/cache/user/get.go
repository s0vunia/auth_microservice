package user

import (
	"context"
	"errors"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/s0vunia/auth_microservice/internal/cache/user/converter"
	modelCache "github.com/s0vunia/auth_microservice/internal/cache/user/model"
	"github.com/s0vunia/auth_microservice/internal/model"
	"strconv"
)

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
