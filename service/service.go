package service

import "password-lock/repository"

type Service struct {
	userRepository repository.UserRepository
}

func NewService(userRepo repository.UserRepository) Service {
	return Service{
		userRepository: userRepo,
	}
}
