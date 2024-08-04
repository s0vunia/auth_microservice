package user

import (
	"context"
	"strconv"
)

func (c *cacheImplementation) Delete(ctx context.Context, id int64) error {
	idStr := strconv.FormatInt(id, 10)
	err := c.cl.Del(ctx, idStr)
	if err != nil {
		return err
	}
	return nil
}
