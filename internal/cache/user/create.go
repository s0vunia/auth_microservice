package user

import (
	"context"
	"strconv"

	"github.com/s0vunia/auth_microservice/internal/cache/user/converter"
	"github.com/s0vunia/auth_microservice/internal/model"
)

func (c *cacheImplementation) Create(ctx context.Context, user *model.User) error {
	userCache := converter.ToUserCacheFromModel(user)
	idStr := strconv.FormatInt(userCache.ID, 10)
	err := c.cl.HashSet(ctx, idStr, userCache)
	if err != nil {
		return err
	}
	return nil
}
