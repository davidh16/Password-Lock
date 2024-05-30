package service

import (
	"context"
	"github.com/google/uuid"
	"time"
)

func (s Service) GenerateAndSaveSessionKey(userUuid string, sessionLifeTime time.Duration) (string, error) {
	ctx := context.Background()

	sessionKey := uuid.New().String()

	err := s.redis.Set(ctx, sessionKey, userUuid, sessionLifeTime).Err()
	if err != nil {
		return "", err
	}

	return sessionKey, nil
}

func (s Service) TerminateSession(sessionId string) {
	ctx := context.Background()
	_ = s.redis.Del(ctx, sessionId)
}
