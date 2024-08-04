package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/s0vunia/auth_microservice/internal/api/user"
	"github.com/s0vunia/auth_microservice/internal/service"
	serviceMocks "github.com/s0vunia/auth_microservice/internal/service/mocks"
	desc "github.com/s0vunia/auth_microservice/pkg/auth_v1"
	"github.com/stretchr/testify/require"
)

func TestImplementation_IsExists(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.IsExistsRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		ids = []int64{1, 2, 3, 4, 5}

		serviceErr = fmt.Errorf("service error")

		req = &desc.IsExistsRequest{
			Ids: ids,
		}

		res1 = &desc.IsExistsResponse{
			Exists: true,
		}
		res2 = &desc.IsExistsResponse{
			Exists: false,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.IsExistsResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case - true",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res1,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.IsExistsMock.Expect(ctx, ids).Return(true, nil)
				return mock
			},
		},
		{
			name: "success case - false",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res2,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.IsExistsMock.Expect(ctx, ids).Return(false, nil)
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
				mock.IsExistsMock.Expect(ctx, ids).Return(false, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			serviceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(serviceMock)

			newID, err := api.IsExists(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
