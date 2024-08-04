package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/s0vunia/auth_microservice/internal/cache"
	mocks2 "github.com/s0vunia/auth_microservice/internal/cache/mocks"
	"github.com/s0vunia/auth_microservice/internal/model"
	"github.com/s0vunia/auth_microservice/internal/repository"
	repoMocks "github.com/s0vunia/auth_microservice/internal/repository/mocks"
	userService "github.com/s0vunia/auth_microservice/internal/service/user"
	"github.com/s0vunia/platform_common/pkg/db"
	"github.com/s0vunia/platform_common/pkg/db/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type logRepositoryMockFunc func(mc *minimock.Controller) repository.LogRepository
	type cacheMockFunc func(mc *minimock.Controller) cache.UserCache
	type txManagerMockFunc func(f func(context.Context) error, mc *minimock.Controller) db.TxManager

	type args struct {
		ctx         context.Context
		userRepoReq *model.UserCreate
		logRepoReq  *model.LogCreate
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Animal()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, false, 10)
		role     = gofakeit.IntRange(1, 2)

		repoErr = fmt.Errorf("repo error")

		userRepoReq = &model.UserCreate{
			Name:     name,
			Email:    email,
			Role:     model.Role(role),
			Password: password,
		}

		idLog      = gofakeit.Int64()
		message    = gofakeit.Animal()
		logRepoReq = &model.LogCreate{
			Message: message,
		}

		cacheUser = &model.User{
			ID: id,
			Info: model.UserInfo{
				Name:  userRepoReq.Name,
				Email: userRepoReq.Email,
				Role:  userRepoReq.Role,
			},
		}
	)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
		logRepositoryMock  logRepositoryMockFunc
		cacheMock          cacheMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:         ctx,
				userRepoReq: userRepoReq,
				logRepoReq:  logRepoReq,
			},
			want: id,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, userRepoReq).Return(id, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, logRepoReq).Return(idLog, nil)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := mocks2.NewUserCacheMock(mc)
				mock.CreateMock.Expect(ctx, cacheUser).Return(nil)
				return mock
			},
			txManagerMock: func(f func(context.Context) error, mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Expect(ctx, f).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx:         ctx,
				userRepoReq: userRepoReq,
				logRepoReq:  logRepoReq,
			},
			want: 0,
			err:  repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, userRepoReq).Return(0, repoErr)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, logRepoReq).Return(idLog, nil)
				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := mocks2.NewUserCacheMock(mc)
				mock.CreateMock.Expect(ctx, cacheUser).Return(nil)
				return mock
			},
			txManagerMock: func(f func(context.Context) error, mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Expect(ctx, f).Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepositoryMock(mc)
			logRepoMock := tt.logRepositoryMock(mc)
			cacheMock := tt.cacheMock(mc)
			txManagerMock := tt.txManagerMock(func(ctx context.Context) error {
				var errTx error
				id, errTx = userRepoMock.Create(ctx, userRepoReq)
				if errTx != nil {
					return errTx
				}

				errTx = cacheMock.Create(ctx, cacheUser)

				if errTx != nil {
					return errTx
				}

				_, errTx = logRepoMock.Create(ctx, logRepoReq)

				if errTx != nil {
					return errTx
				}

				return nil
			}, mc)

			service := userService.NewService(
				userRepoMock,
				logRepoMock,
				cacheMock,
				txManagerMock,
			)

			newID, err := service.Create(tt.args.ctx, tt.args.userRepoReq)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
