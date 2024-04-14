package repository

import (
	"gorm.io/gorm"
	"password-lock/models"
)

type TokenRepository interface {
	Db() *gorm.DB
	FindTokenByToken(token string) (*models.Token, error)
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *tokenRepository {
	return &tokenRepository{
		db: db,
	}
}

func (r tokenRepository) Db() *gorm.DB {
	return r.db
}

func (r tokenRepository) FindTokenByToken(token string) (*models.Token, error) {
	var tokenModel models.Token
	result := r.db.Where("token=?", token).First(&tokenModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tokenModel, nil
}
