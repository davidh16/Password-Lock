package service

import (
	"github.com/redis/go-redis/v9"
	"password-lock/repository"
)

type Service struct {
	redis            *redis.Client
	userRepository   repository.UserRepository
	entityRepository repository.EntityRepository
}

/* */
func NewService(redis *redis.Client,
	userRepo repository.UserRepository,
	entityRepo repository.EntityRepository,
) Service {

	return Service{
		redis:            redis,
		userRepository:   userRepo,
		entityRepository: entityRepo,
	}

}
