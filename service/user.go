package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"password-lock/models"
)

func (s Service) RegisterUser(ctx *gin.Context, user *models.User) (*models.User, error) {

	tx := s.userRepository.Db().Begin()
	err := setTransaction(ctx, []*gorm.DB{tx})
	if err != nil {
		return nil, err
	}

	result := tx.Table("users").Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s Service) VerifyUser(ctx *gin.Context, userUuid string, password string) (*models.User, error) {

	tx := s.userRepository.Db().Begin()
	err := setTransaction(ctx, []*gorm.DB{tx})
	if err != nil {
		return nil, err
	}

	var user models.User
	result := tx.Where("uuid=? AND active = FALSE", userUuid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	user.Active = true
	user.Password = password

	result = tx.Table("users").Where("uuid=? AND active = FALSE", userUuid).Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (s Service) Authenticate(credentials models.User) (*models.User, error) {
	user, err := s.userRepository.FindUserByEmailAddress(credentials.EmailAddress)
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

func (s Service) IfEmailAddressExists(emailAddress string) (error, bool) {
	user, err := s.userRepository.FindUserByEmailAddress(emailAddress)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false
		} else {
			return err, false
		}
	}
	if user != nil {
		return nil, true
	} else {
		return nil, false
	}
}
