package user

import (
	"context"
	"fmt"

	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/model"
)

func (s serv) Create(ctx context.Context, userCreate *model.UserCreate) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
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
