package user

import "context"

func (s serv) IsExists(ctx context.Context, ids []int64) (bool, error) {
	return s.userRepository.IsExists(ctx, ids)
}
