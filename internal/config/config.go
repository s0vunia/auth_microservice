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

// JWTConfig - config for JWT
type JWTConfig interface {
	AccessSecretKey() string
	AccessExpiration() time.Duration
	RefreshSecretKey() string
	RefreshExpiration() time.Duration
	AuthPrefix() string
}

// SwaggerConfig - config for swagger
type SwaggerConfig interface {
	Address() string
}

// PrometheusConfig - config for prometheus
type PrometheusConfig interface {
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

// LoggerConfig - config for logger
type LoggerConfig interface {
	Level() string
	FileName() string
	MaxSize() int
	MaxAge() int
	MaxBackups() int
}

// Load - loads config from .env
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
