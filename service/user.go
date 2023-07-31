package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"password-lock/models"
)

func (s Service) RegisterUser(user models.User) (*models.User, error) {
	result := s.userRepository.Db().Table("users").Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (s Service) Authenticate(credentials models.User) (*models.User, error) {
	user, err := s.userRepository.FindUserByEmailAddress(credentials)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user does not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
