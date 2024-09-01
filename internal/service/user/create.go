package user

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/model"
	"github.com/s0vunia/auth_microservice/internal/utils"
)

func (s serv) Create(ctx context.Context, userCreate *model.UserCreate) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		password, errTx := utils.HashPassword(userCreate.Password)
		if errTx != nil {
			return errTx
		}
		userCreate.Password = password

		id, errTx = s.userRepository.Create(ctx, userCreate)
		if errTx != nil {
			return errTx
		}

		errTx = s.cache.Create(ctx, &model.User{
			ID: id,
			Info: model.UserInfo{
				Name:  userCreate.Name,
				Email: userCreate.Email,
				Role:  userCreate.Role,
			},
		})

		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}
	return id, nil
}
