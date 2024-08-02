package converter

import (
	"github.com/s0vunia/auth_microservice/internal/model"
	modelRepo "github.com/s0vunia/auth_microservice/internal/repository/user/model"
)

// ToUserFromRepo converts a User object from the model package to a User object from the repository package.
func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Info:      ToUserInfoFromRepo(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserInfoFromRepo converts a UserInfo object from the model package to a UserInfo object from the repository package.
func ToUserInfoFromRepo(info modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:  info.Name,
		Email: info.Email,
		Role:  model.Role(info.Role),
	}
}

// ToUserCreateFromService converts a UserCreate object from the model package to a UserCreate object from the repository package.
func ToUserCreateFromService(userCreate *model.UserCreate) *modelRepo.UserCreate {
	return &modelRepo.UserCreate{
		Name:         userCreate.Name,
		Email:        userCreate.Email,
		Role:         modelRepo.Role(userCreate.Role),
		PasswordHash: []byte(userCreate.Password),
	}
}
