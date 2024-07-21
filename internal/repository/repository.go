package repository

import (
	"context"

	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/model"
)

// UserRepository represents a repository for user entities.
type UserRepository interface {
	Create(ctx context.Context, userCreate *model.UserCreate) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, id int64, userUpdate *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
}
