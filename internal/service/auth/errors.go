package auth

import "errors"

var (
	// ErrInvalidRefreshToken - invalid refresh token
	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	// ErrInvalidAccessToken - invalid access token
	ErrInvalidAccessToken = errors.New("invalid access token")

	// ErrWrongPassword - wrong password
	ErrWrongPassword = errors.New("wrong password")

	// ErrGenerateToken - failed to generate token
	ErrGenerateToken = errors.New("failed to generate token")
)
