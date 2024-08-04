package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/s0vunia/auth_microservice/internal/api/user"
	"github.com/s0vunia/auth_microservice/internal/model"
	"github.com/s0vunia/auth_microservice/internal/service"
	serviceMocks "github.com/s0vunia/auth_microservice/internal/service/mocks"
	desc "github.com/s0vunia/auth_microservice/pkg/auth_v1"
	"github.com/stretchr/testify/require"
)

func TestImplementation_Create(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Animal()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, false, 10)
		role     = gofakeit.IntRange(1, 2)

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			UserCreate: &desc.UserCreate{
				Info: &desc.UserInfo{
					Name:  name,
					Email: email,
					Role:  desc.Role(role),
				},
				Password:        password,
				PasswordConfirm: password,
			},
		}
		reqWithDiffPasswords = &desc.CreateRequest{
			UserCreate: &desc.UserCreate{
				Info: &desc.UserInfo{
					Name:  name,
					Email: email,
					Role:  desc.Role(role),
				},
				Password:        password,
				PasswordConfirm: password + "a",
			},
		}

		userCreate = model.UserCreate{
			Name:     name,
			Email:    email,
			Role:     model.Role(role),
			Password: password,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, &userCreate).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, &userCreate).Return(0, serviceErr)
				return mock
			},
		},
		{
			name: "password and password confirm do not match",
			args: args{
				ctx: ctx,
				req: reqWithDiffPasswords,
			},
			want: nil,
			err:  user.PasswordConfirmErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			createServiceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(createServiceMock)

			newID, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
