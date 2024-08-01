package app

import (
	"context"
	"log"

	cacheCl "github.com/s0vunia/platform_common/pkg/cache"
	"github.com/s0vunia/platform_common/pkg/cache/redis"
	"github.com/s0vunia/platform_common/pkg/closer"
	"github.com/s0vunia/platform_common/pkg/db"
	"github.com/s0vunia/platform_common/pkg/db/pg"
	"github.com/s0vunia/platform_common/pkg/db/transaction"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/cache"
	userCache "github.com/s0vunia/auth_microservices_course_boilerplate/internal/cache/user"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/config/env"

	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/api/user"

	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/config"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/repository"
	logRepository "github.com/s0vunia/auth_microservices_course_boilerplate/internal/repository/log"
	userRepository "github.com/s0vunia/auth_microservices_course_boilerplate/internal/repository/user"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/service"
	userService "github.com/s0vunia/auth_microservices_course_boilerplate/internal/service/user"
)

type serviceProvider struct {
	pgConfig    config.PGConfig
	grpcConfig  config.GRPCConfig
	redisConfig config.RedisConfig

	dbClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cacheCl.Client

	userRepository repository.UserRepository
	logsRepository repository.LogRepository
	cache          cache.UserCache

	userService service.UserService

	userImpl *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := env.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
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

func (s *serviceProvider) LogsRepository(ctx context.Context) repository.LogRepository {
	if s.logsRepository == nil {
		s.logsRepository = logRepository.NewRepository(s.DBClient(ctx))
	}

	return s.logsRepository
}

func (s *serviceProvider) Cache() cache.UserCache {
	if s.cache == nil {
		s.cache = userCache.NewCache(s.RedisClient())
	}
	return s.cache
}

func (s *serviceProvider) NoteService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.LogsRepository(ctx),
			s.Cache(),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.NoteService(ctx))
	}

	return s.userImpl
}
