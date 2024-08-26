package user

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/model"
)

func (s serv) Get(ctx context.Context, id int64) (*model.User, error) {
	var user *model.User

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		user, errTx = s.cache.Get(ctx, id)
		if errTx != nil {
			user, errTx = s.userRepository.Get(ctx, id)
			if errTx != nil {
				return errTx
			}

			errTx = s.cache.Create(ctx, user)
			if errTx != nil {
				return errTx
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}
