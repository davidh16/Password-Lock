package service

import (
	"context"
	"github.com/google/uuid"
	"password-lock/models"
	"time"
)

func (s Service) GenerateAndSaveSessionKey(user models.User) (string, error) {
	ctx := context.Background()

	sessionKey := uuid.New().String()

	err := s.redis.Set(ctx, sessionKey, user.EmailAddress, time.Minute*10).Err()
	if err != nil {
		return "", err
	}

	return sessionKey, nil
}
