package auth

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/model"
	"github.com/s0vunia/auth_microservice/internal/utils"
)

func (s serv) Login(ctx context.Context, login *model.UserLogin) (string, error) {
	var refreshToken string
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		passHash, errTx := s.userRepository.GetPassword(ctx, login.ID)
		if errTx != nil {
			return errTx
		}

		if !utils.VerifyPassword(string(passHash), login.Password) {
			return ErrWrongPassword
		}
		user, errTx := s.userRepository.Get(ctx, login.ID)
		if errTx != nil {
			return errTx
		}

		refreshToken, errTx = utils.GenerateToken(*user,
			[]byte(s.refreshTokenSecretKey),
			s.refreshTokenExpiration,
		)
		if errTx != nil {
			return ErrGenerateToken
		}

		return nil
	})
	return refreshToken, err
}
