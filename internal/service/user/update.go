package user

import (
	"context"
	"fmt"

	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/model"
)

func (s serv) Update(ctx context.Context, id int64, userUpdate *model.UserUpdate) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.userRepository.Update(ctx, id, userUpdate)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.logsRepository.Create(ctx, &model.LogCreate{
			Message: fmt.Sprintf("User %d updated", id),
		})

		if errTx != nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
