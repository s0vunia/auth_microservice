package auth

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/model"
	"github.com/s0vunia/auth_microservice/internal/utils"
)

func (s serv) GetAccessToken(_ context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(s.refreshTokenSecretKey))
	if err != nil {
		return "", ErrInvalidAccessToken
	}

	refreshToken, err = utils.GenerateToken(model.User{
		ID: claims.ID,
		Info: model.UserInfo{
			Role: claims.Role,
		},
	},
		[]byte(s.accessTokenSecretKey),
		s.accessTokenExpiration,
	)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
