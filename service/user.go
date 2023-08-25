package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"password-lock/models"
)

func (s Service) RegisterUser(ctx *gin.Context, user models.User) (*models.User, error) {
	tx := s.entityRepository.Db().Begin()
	ctx.Set("tx", tx)

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

func (s Service) Authorize(userUuid string, password string) error {
	user, err := s.userRepository.FindUserByUuid(userUuid)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user does not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return err
	}

	return err
}

func (s Service) Me(ctx *gin.Context) string {
	sessionUuid, err := ctx.Cookie("session")
	if err != nil {
		return ""
	}

	loggedInUser, err := s.redis.Get(context.Background(), sessionUuid).Result()
	if err != nil {
		return ""
	}

	return loggedInUser
}
