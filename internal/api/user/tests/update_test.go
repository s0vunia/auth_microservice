package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/s0vunia/auth_microservice/internal/api/user"
	"github.com/s0vunia/auth_microservice/internal/model"
	"github.com/s0vunia/auth_microservice/internal/service"
	serviceMocks "github.com/s0vunia/auth_microservice/internal/service/mocks"
	desc "github.com/s0vunia/auth_microservice/pkg/auth_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

func TestImplementation_Update(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Animal()
		email = gofakeit.Email()
		role  = gofakeit.IntRange(1, 2)

		serviceErr = fmt.Errorf("service error")

		req = &desc.UpdateRequest{
			Id: id,
			UserUpdate: &desc.UserUpdate{
				Name:  wrapperspb.String(name),
				Email: wrapperspb.String(email),
				Role:  desc.Role(role),
			},
		}
		userUpdate = model.UserUpdate{
			Name:  &name,
			Email: &email,
			Role:  model.Role(role),
		}

		res = &emptypb.Empty{}
	)
	defer t.Cleanup(mc.Finish)
	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
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
				mock.UpdateMock.Expect(ctx, id, &userUpdate).Return(nil)
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
				mock.UpdateMock.Expect(ctx, id, &userUpdate).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			updateServiceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(updateServiceMock)

			newID, err := api.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
