package user

import (
	"context"
	"fmt"

	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/model"
)

func (s serv) Get(ctx context.Context, id int64) (*model.User, error) {
	var user *model.User

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		user, errTx = s.cache.Get(ctx, id)
		message := fmt.Sprintf("User %d got from cache", id)
		if errTx != nil {
			message = fmt.Sprintf("User %d got from db", id)
			user, errTx = s.userRepository.Get(ctx, id)
			if errTx != nil {
				return errTx
			}

			errTx = s.cache.Create(ctx, user)
			if errTx != nil {
				return errTx
			}
		}

		_, errTx = s.logsRepository.Create(ctx, &model.LogCreate{
			Message: message,
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}
