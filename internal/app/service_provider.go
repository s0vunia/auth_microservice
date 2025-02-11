package app

import (
	"context"

	"github.com/IBM/sarama"
	accessImplem "github.com/s0vunia/auth_microservice/internal/api/access"
	"github.com/s0vunia/auth_microservice/internal/api/auth"
	"github.com/s0vunia/auth_microservice/internal/logger"
	userSaverConsumer "github.com/s0vunia/auth_microservice/internal/service/consumer/user_saver"
	cacheCl "github.com/s0vunia/platform_common/pkg/cache"
	"github.com/s0vunia/platform_common/pkg/cache/redis"
	"github.com/s0vunia/platform_common/pkg/closer"
	"github.com/s0vunia/platform_common/pkg/db"
	"github.com/s0vunia/platform_common/pkg/db/pg"
	"github.com/s0vunia/platform_common/pkg/db/transaction"
	"github.com/s0vunia/platform_common/pkg/kafka"
	"go.uber.org/zap"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/s0vunia/auth_microservice/internal/cache"
	userCache "github.com/s0vunia/auth_microservice/internal/cache/user"
	"github.com/s0vunia/auth_microservice/internal/config/env"
	kafkaConsumer "github.com/s0vunia/platform_common/pkg/kafka/consumer"

	"github.com/s0vunia/auth_microservice/internal/api/user"

	"github.com/s0vunia/auth_microservice/internal/config"
	"github.com/s0vunia/auth_microservice/internal/repository"
	userRepository "github.com/s0vunia/auth_microservice/internal/repository/user"
	"github.com/s0vunia/auth_microservice/internal/service"
	accessService "github.com/s0vunia/auth_microservice/internal/service/access"
	authService "github.com/s0vunia/auth_microservice/internal/service/auth"
	userService "github.com/s0vunia/auth_microservice/internal/service/user"
)

type serviceProvider struct {
	pgConfig            config.PGConfig
	grpcConfig          config.GRPCConfig
	httpConfig          config.HTTPConfig
	jwtConfig           config.JWTConfig
	swaggerConfig       config.SwaggerConfig
	prometheusConfig    config.SwaggerConfig
	redisConfig         config.RedisConfig
	kafkaConsumerConfig config.KafkaConsumerConfig
	loggerConfig        config.LoggerConfig

	dbClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cacheCl.Client

	userRepository repository.UserRepository
	cache          cache.UserCache

	userService   service.UserService
	authService   service.AuthService
	accessService service.AccessService

	userSaverConsumer service.ConsumerService

	userImpl   *user.Implementation
	authImpl   *auth.Implementation
	accessImpl *accessImplem.Implementation

	consumer             kafka.Consumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *kafkaConsumer.GroupHandler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			logger.Fatal(
				"failed to get pg config",
				zap.Error(err),
			)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			logger.Fatal(
				"failed to get grpc config",
				zap.Error(err),
			)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			logger.Fatal(
				"failed to get http config",
				zap.Error(err),
			)
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := env.NewJWTConfig()
		if err != nil {
			logger.Fatal(
				"failed to get jwt config",
				zap.Error(err),
			)
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			logger.Fatal(
				"failed to get swagger config",
				zap.Error(err),
			)
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) PrometheusConfig() config.PrometheusConfig {
	if s.prometheusConfig == nil {
		cfg, err := env.NewPrometheusConfig()
		if err != nil {
			logger.Fatal(
				"failed to get prometheus config",
				zap.Error(err),
			)
		}

		s.prometheusConfig = cfg
	}

	return s.prometheusConfig
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := env.NewRedisConfig()
		if err != nil {
			logger.Fatal(
				"failed to get redis config",
				zap.Error(err),
			)
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) KafkaConsumerConfig() config.KafkaConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := env.NewKafkaConsumerConfig()
		if err != nil {
			logger.Fatal(
				"failed to get kafka consumer config",
				zap.Error(err),
			)
		}

		s.kafkaConsumerConfig = cfg
	}

	return s.kafkaConsumerConfig
}

func (s *serviceProvider) LoggerConfig() config.LoggerConfig {
	if s.loggerConfig == nil {
		cfg, err := env.NewLoggerConfig()
		if err != nil {
			logger.Fatal(
				"failed to get logger config",
				zap.Error(err),
			)
		}
		s.loggerConfig = cfg
	}
	return s.loggerConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			logger.Fatal(
				"failed to get db client",
				zap.Error(err),
			)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			logger.Fatal(
				"failed to ping db",
				zap.Error(err),
			)
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

func (s *serviceProvider) RedisClient() cacheCl.Client {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig().ConnectionTimeout())
	}

	return s.redisClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) Cache() cache.UserCache {
	if s.cache == nil {
		s.cache = userCache.NewCache(s.RedisClient())
	}
	return s.cache
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.Cache(),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.UserRepository(ctx),
			s.Cache(),
			s.TxManager(ctx),
			s.JWTConfig().RefreshSecretKey(),
			s.JWTConfig().RefreshExpiration(),
			s.JWTConfig().AccessSecretKey(),
			s.JWTConfig().AccessExpiration(),
		)
	}

	return s.authService
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewService(
			s.Cache(),
			s.TxManager(ctx),
			s.JWTConfig().AuthPrefix(),
			s.JWTConfig().AccessSecretKey(),
		)
	}

	return s.accessService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}

func (s *serviceProvider) AccessImpl(ctx context.Context) *accessImplem.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = accessImplem.NewImplementation(s.AccessService(ctx))
	}

	return s.accessImpl
}

func (s *serviceProvider) UserSaverConsumer(ctx context.Context) service.ConsumerService {
	if s.userSaverConsumer == nil {
		s.userSaverConsumer = userSaverConsumer.NewService(
			s.UserRepository(ctx),
			s.Consumer(),
		)
	}

	return s.userSaverConsumer
}

func (s *serviceProvider) Consumer() kafka.Consumer {
	if s.consumer == nil {
		s.consumer = kafkaConsumer.NewConsumer(
			s.ConsumerGroup(),
			s.ConsumerGroupHandler(),
		)
		closer.Add(s.consumer.Close)
	}

	return s.consumer
}

func (s *serviceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			s.KafkaConsumerConfig().Brokers(),
			s.KafkaConsumerConfig().GroupID(),
			s.KafkaConsumerConfig().Config(),
		)
		if err != nil {
			logger.Fatal(
				"failed to create consumer group",
				zap.Error(err),
			)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

func (s *serviceProvider) ConsumerGroupHandler() *kafkaConsumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = kafkaConsumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}
