package service

import (
	"github.com/redis/go-redis/v9"
	"password-lock/repository"
)

type Service struct {
	redis          *redis.Client
	userRepository repository.UserRepository
}

func NewService(redis *redis.Client, userRepo repository.UserRepository) Service {
	return Service{
		redis:          redis,
		userRepository: userRepo,
	}
}
