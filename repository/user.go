package repository

import (
	"gorm.io/gorm"
	"password-lock/db"
	"password-lock/models"
)

type UserRepository interface {
	Db() *gorm.DB
	FindUserByUuid(userUuid string) (*models.User, error)
	FindAllSecurityQuestions() ([]models.PersonalQuestion, error)
	FindUserByEmailAddress(emailAddress string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r userRepository) Db() *gorm.DB {
	return r.db
}

func (r userRepository) FindUserByEmailAddress(emailAddress string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email_address = ?", emailAddress).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

//func (r userRepository) FindUserByEmailAddress(emailAddress string) (*models.User, error) {
//	var user models.User
//	result := r.db.Where("email_address = ? AND active = TRUE", emailAddress).First(&user)
//	if result.Error != nil {
//		return nil, result.Error
//	}
//	return &user, nil
//}
//func (r userRepository) FindUnverifiedUserByEmailAddress(emailAddress string) (*models.User, error) {
//	var user models.User
//	result := r.db.Where("email_address = ? AND active = false", emailAddress).First(&user)
//	if result.Error != nil {
//		return nil, result.Error
//	}
//	return &user, nil
//}

func (r userRepository) FindUserByUuid(userUuid string) (*models.User, error) {
	var foundUser models.User
	result := r.db.Where("uuid = ?", userUuid).First(&foundUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}

func (r userRepository) FindAllSecurityQuestions() ([]models.PersonalQuestion, error) {
	var securityQuestions []models.PersonalQuestion
	result := r.Db().Table(db.PERSONAL_QUESTIONS_TABLE).Find(&securityQuestions)
	if result.Error != nil {
		return nil, result.Error
	}
	return securityQuestions, nil
}
