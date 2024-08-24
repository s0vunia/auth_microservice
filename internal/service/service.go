package service

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/model"
)

// UserService represents a service for user entities.
type UserService interface {
	Create(ctx context.Context, userCreate *model.UserCreate) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, id int64, userUpdate *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
	IsExists(ctx context.Context, ids []int64) (bool, error)
}

// ConsumerService represents a service for consumer entities.
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

// AuthService represents a service for auth entities.
type AuthService interface {
	Login(ctx context.Context, userLogin *model.UserLogin) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

// AccessService represents a service for access entities.
type AccessService interface {
	Check(ctx context.Context, endpointAddress string) error
}
