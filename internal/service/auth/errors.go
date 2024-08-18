package auth

import "errors"

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrWrongPassword       = errors.New("wrong password")
	ErrGenerateToken       = errors.New("failed to generate token")
)
