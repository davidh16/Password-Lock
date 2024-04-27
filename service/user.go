package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"password-lock/models"
)

func (s Service) RegisterUser(ctx *gin.Context, user *models.User) (*models.User, error) {

	result := s.userRepository.Db().Table("users").Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s Service) VerifyUser(ctx *gin.Context, userUuid string, password string) (*models.User, error) {

	var user models.User
	result := s.userRepository.Db().Where("uuid=? AND active = FALSE", userUuid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	user.Active = true
	user.Password = password

	result = s.userRepository.Db().Set("encrypt-password", true).Table("users").Where("uuid=? AND active = FALSE", userUuid).Save(&user)
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

	result := tx.Table("users").Where("uuid=? AND active = TRUE AND completed = FALSE", user.Uuid).Save(user)
	if result.Error != nil {
		return nil, result.Error
	}

	result = tx.Table("user_personal_questions").Create(personalQuestions)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
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

func (s Service) Me(ctx *gin.Context) (*models.User, error) {

	//me, err := s.userRepository.FindUserByUuid(ctx.Value("me").(string))
	//if err != nil {
	//	return nil, err
	//}

	me, err := s.userRepository.FindUserByUuid("3174f700-389a-4cf0-8636-8c520145baf5")
	if err != nil {
		return nil, err
	}

	return me, nil
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

func (s Service) GetUserByEmailAddress(emailAddress string) (*models.User, error) {
	user, err := s.userRepository.FindUserByEmailAddress(emailAddress)
	if err != nil {
		return nil, err
	}
	return user, nil
}
