package user

import (
	"context"
	"fmt"

	"github.com/s0vunia/auth_microservice/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (s serv) Create(ctx context.Context, userCreate *model.UserCreate) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		password, errTx := bcrypt.GenerateFromPassword([]byte(userCreate.Password), 10)
		if errTx != nil {
			return errTx
		}
		userCreate.Password = string(password)

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

		_, errTx = s.logsRepository.Create(ctx, &model.LogCreate{
			Message: fmt.Sprintf("User %d created", id),
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
