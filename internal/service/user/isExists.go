package user

import "context"

func (s serv) IsExists(ctx context.Context, id int64) (bool, error) {
	return s.userRepository.IsExists(ctx, id)
}
