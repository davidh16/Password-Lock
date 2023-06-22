package service

import (
	"password-lock/models"
)

func (s Service) RegisterUser(user models.User) (*models.User, error) {
	result := s.repository.Db().Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
