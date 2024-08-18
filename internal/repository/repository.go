package repository

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/model"
)

// UserRepository represents a repository for user entities.
type UserRepository interface {
	Create(ctx context.Context, userCreate *model.UserCreate) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	GetPassword(ctx context.Context, id int64) (string, error)
	Update(ctx context.Context, id int64, userUpdate *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
	IsExists(ctx context.Context, ids []int64) (bool, error)
}

// LogRepository represents a repository for log entities.
type LogRepository interface {
	Create(ctx context.Context, logCreate *model.LogCreate) (int64, error)
}
