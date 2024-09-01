package user_saver

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/s0vunia/auth_microservice/internal/model"
)

func (s *service) UserSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	userInfo := &model.UserCreate{}
	err := json.Unmarshal(msg.Value, userInfo)
	if err != nil {
		return err
	}

	_, err = s.userRepository.Create(ctx, userInfo)
	if err != nil {
		return err
	}

	return nil
}
