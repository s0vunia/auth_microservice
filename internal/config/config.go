package config

import (
	"time"

	"github.com/IBM/sarama"
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

// HTTPConfig - config for HTTP server
type HTTPConfig interface {
	Address() string
	ReadHeaderTimeout() time.Duration
}

// SwaggerConfig - config for swagger
type SwaggerConfig interface {
	Address() string
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

// KafkaConsumerConfig - config for kafka consumer
type KafkaConsumerConfig interface {
	Brokers() []string
	GroupID() string
	Config() *sarama.Config
}

// Load - loads config from .env
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
