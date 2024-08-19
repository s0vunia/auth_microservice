package user_saver

import (
	"context"

	"github.com/s0vunia/auth_microservice/internal/repository"
	def "github.com/s0vunia/auth_microservice/internal/service"
	"github.com/s0vunia/platform_common/pkg/kafka"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
	consumer       kafka.Consumer
}

// NewService creates a new user saver service.
func NewService(
	userRepository repository.UserRepository,
	consumer kafka.Consumer,
) *service {
	return &service{
		userRepository: userRepository,
		consumer:       consumer,
	}
}

// RunConsumer runs user saver consumer.
func (s *service) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

func (s *service) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- s.consumer.Consume(ctx, "test-topic", s.UserSaveHandler)
	}()

	return errChan
}
