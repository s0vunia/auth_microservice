package env

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/s0vunia/auth_microservice/internal/config"
)

const (
	accessSecretKeyEnvName   = "ACCESS_TOKEN_SECRET_KEY"
	accessExpirationEnvName  = "ACCESS_TOKEN_EXPIRATION_SEC"
	refreshSecretKeyEnvName  = "REFRESH_TOKEN_SECRET_KEY"
	refreshExpirationEnvName = "REFRESH_TOKEN_EXPIRATION_SEC"
	authPrefixEnvName        = "AUTH_PREFIX"
)

type jwtConfig struct {
	accessSecretKey   string
	accessExpiration  time.Duration
	refreshSecretKey  string
	refreshExpiration time.Duration
	authPrefix        string
}

// NewJWTConfig - creates new jwt config
func NewJWTConfig() (config.JWTConfig, error) {
	accessSecretKey := os.Getenv(accessSecretKeyEnvName)
	if len(accessSecretKey) == 0 {
		return nil, errors.New("access secret key not found")
	}

	accessExpiration := os.Getenv(accessExpirationEnvName)
	if len(accessExpiration) == 0 {
		return nil, errors.New("access expiration not found")
	}

	accessDurSec, err := strconv.Atoi(accessExpiration)
	if err != nil {
		return nil, err
	}
	accessExpirationDuration := time.Second * time.Duration(accessDurSec)

	refreshSecretKey := os.Getenv(refreshSecretKeyEnvName)
	if len(refreshSecretKey) == 0 {
		return nil, errors.New("refresh secret key not found")
	}

	refreshExpiration := os.Getenv(refreshExpirationEnvName)
	if len(refreshExpiration) == 0 {
		return nil, errors.New("refresh expiration not found")
	}

	refreshDurSec, err := strconv.Atoi(refreshExpiration)
	if err != nil {
		return nil, err
	}
	refreshExpirationDuration := time.Second * time.Duration(refreshDurSec)

	authPrefix := os.Getenv(authPrefixEnvName) + " "
	if len(authPrefix) == 0 {
		return nil, errors.New("auth prefix not found")
	}

	return &jwtConfig{
		accessSecretKey:   accessSecretKey,
		accessExpiration:  accessExpirationDuration,
		refreshSecretKey:  refreshSecretKey,
		refreshExpiration: refreshExpirationDuration,
		authPrefix:        authPrefix,
	}, nil
}

// AccessSecretKey - returns access secret key
func (cfg *jwtConfig) AccessSecretKey() string {
	return cfg.accessSecretKey
}

// AccessExpiration - returns access expiration
func (cfg *jwtConfig) AccessExpiration() time.Duration {
	return cfg.accessExpiration
}

// RefreshSecretKey - returns refresh secret key
func (cfg *jwtConfig) RefreshSecretKey() string {
	return cfg.refreshSecretKey
}

// RefreshExpiration - returns refresh expiration
func (cfg *jwtConfig) RefreshExpiration() time.Duration {
	return cfg.refreshExpiration
}

// AuthPrefix - returns auth prefix
func (cfg *jwtConfig) AuthPrefix() string {
	return cfg.authPrefix
}
