package user

import (
	"context"
)

func (s serv) IsExists(ctx context.Context, ids []int64) (bool, error) {
	var isExists bool

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		isExists, errTx = s.userRepository.IsExists(ctx, ids)
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
