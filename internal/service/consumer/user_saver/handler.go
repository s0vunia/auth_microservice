package user_saver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/s0vunia/auth_microservice/internal/model"
)

func (s *service) UserSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	userInfo := &model.UserCreate{}
	err := json.Unmarshal(msg.Value, userInfo)
	if err != nil {
		return err
	}

	id, err := s.userRepository.Create(ctx, userInfo)
	if err != nil {
		return err
	}

	log.Printf("user with id %d created\n", id)

	return nil
}
