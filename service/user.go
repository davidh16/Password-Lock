package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"password-lock/db"
	"password-lock/models"
)

func (s Service) RegisterUser(ctx *gin.Context, user *models.User) (*models.User, error) {

	result := s.userRepository.Db().Table("users").Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s Service) UpdatePassword(ctx *gin.Context, user *models.User, newPassword string) error {
	user.Password = newPassword
	result := s.userRepository.Db().Set("encrypt-password", true).Table(db.USERS_TABLE).Where("uuid=?", user.Uuid).Omit(clause.Associations).Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s Service) VerifyUser(ctx *gin.Context, userUuid string, password string) (*models.User, error) {

	var user models.User
	result := s.userRepository.Db().Table(db.USERS_TABLE).Where("uuid=? AND active = FALSE", userUuid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	user.Active = true
	user.Password = password

	result = s.userRepository.Db().Set("encrypt-password", true).Table(db.USERS_TABLE).Where("uuid=? AND active = FALSE", userUuid).Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (s Service) CompleteRegistration(ctx *gin.Context, user *models.User, personalQuestions []*models.UserPersonalQuestion) (*models.User, error) {

	tx := s.userRepository.Db().Begin()
	err := setTransaction(ctx, []*gorm.DB{tx})
	if err != nil {
		return nil, err
	}

	result := tx.Table(db.USERS_TABLE).Where("uuid=? AND active = TRUE AND completed = FALSE", user.Uuid).Save(user)
	if result.Error != nil {
		return nil, result.Error
	}

	result = tx.Table("user_personal_questions").Create(personalQuestions)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s Service) Authenticate(emailAddress, password string) (*models.User, error) {
	user, err := s.userRepository.FindUserByEmailAddress(emailAddress)
	if err != nil {
		return nil, err
	}

	if user == nil || (user != nil && !user.Active) {
		return nil, errors.New("user does not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
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

func (s Service) Me(ctx *gin.Context) (*models.User, error) {

	fmt.Println("vjv ovdje pukne")

	me, err := s.userRepository.FindUserByUuid(ctx.Value("me").(string))
	if err != nil {
		return nil, err
	}

	return me, nil
}

func (s Service) GetUserByEmailAddress(emailAddress string) (*models.User, error) {
	user, err := s.userRepository.FindUserByEmailAddress(emailAddress)
	if err != nil {
		return nil, err
	}
	return user, nil
}
