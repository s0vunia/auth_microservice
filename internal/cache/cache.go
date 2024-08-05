package cache

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/model"
)

type UserCache interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) error
}
