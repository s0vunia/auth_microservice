package user

import (
	"context"
	"fmt"

	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/model"
)

func (s serv) IsExists(ctx context.Context, ids []int64) (bool, error) {
	var isExists bool

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		isExists, errTx = s.userRepository.IsExists(ctx, ids)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.logsRepository.Create(ctx, &model.LogCreate{
			Message: fmt.Sprintf("Check users exists with ids: %v", ids),
		})

		if errTx != nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	return isExists, nil
}
