package config

import (
	"time"

	"github.com/joho/godotenv"
)

// GRPCConfig - config for GRPC server
type GRPCConfig interface {
	Address() string
}

// PGConfig - config for Postgres
type PGConfig interface {
	DSN() string
}

// RedisConfig - config for Redis
type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

// StorageConfig - config for storage
type StorageConfig interface {
	Mode() string
}

// Load - loads config from .env
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
