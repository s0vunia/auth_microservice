package converter

import (
	"github.com/s0vunia/auth_microservice/internal/model"
	desc "github.com/s0vunia/auth_microservice/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserCreateFromDesc converts a UserCreate object from the desc package to a UserCreate object from the model package.
func ToUserCreateFromDesc(userCreate *desc.UserCreate) *model.UserCreate {
	return &model.UserCreate{
		Name:     userCreate.GetInfo().GetName(),
		Email:    userCreate.GetInfo().GetEmail(),
		Role:     model.Role(userCreate.GetInfo().GetRole()),
		Password: userCreate.GetPassword(),
	}
}

// ToUserFromService converts a User object from the model package to a User object from the desc package.
func ToUserFromService(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}
	return &desc.User{
		Id: user.ID,
		Info: &desc.UserInfo{
			Name:  user.Info.Name,
			Email: user.Info.Email,
			Role:  desc.Role(user.Info.Role),
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// ToUserUpdateFromDesc converts a UserUpdate object from the desc package to a UserUpdate object from the model package.
func ToUserUpdateFromDesc(userUpdate *desc.UserUpdate) *model.UserUpdate {
	var (
		email *string
		name  *string
		role  model.Role
	)

	if userUpdate.GetEmail() != nil {
		email = &userUpdate.GetEmail().Value
	}
	if userUpdate.GetName() != nil {
		name = &userUpdate.GetName().Value
	}
	if userUpdate.GetRole() != desc.Role_UNKNOWN {
		role = (model.Role)(userUpdate.GetRole())
	}
	return &model.UserUpdate{
		Email: email,
		Name:  name,
		Role:  role,
	}
}
