package repository

import (
	"gorm.io/gorm"
	"password-lock/models"
)

type UserRepository interface {
	Db() *gorm.DB
	FindUserByEmailAddress(user models.User) (*models.User, error)
	FindUserByUuid(userUuid string) (*models.User, error)
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

func (r userRepository) FindUserByEmailAddress(user models.User) (*models.User, error) {
	var foundUser models.User
	result := r.db.Where("email_address = ?", user.EmailAddress).First(&foundUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}

func (r userRepository) FindUserByUuid(userUuid string) (*models.User, error) {
	var foundUser models.User
	result := r.db.Where("uuid = ?", userUuid).First(&foundUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}
