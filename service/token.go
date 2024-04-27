package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"password-lock/models"
	"password-lock/utils"
	"time"
)

func (s Service) CreateToken(ctx *gin.Context, userUuid string, tokenType string) (*models.Token, error) {

	token := &models.Token{
		Token:     utils.GenerateToken(),
		TokenType: tokenType,
		UserUuid:  userUuid,
	}

	expireTime := time.Now().Add(models.DefaultTokenExpireTime)
	if tokenType == models.TOKEN_TYPE_FORGOT_PASSWORD {
		expireTime = time.Now().Add(models.DefaultTokenExpireTime)
	}
	if tokenType == models.TOKEN_TYPE_VERIFICATION {
		expireTime = time.Now().Add(time.Hour * 168)
	}

	token.ExpireAt = expireTime

	result := s.tokenRepository.Db().Create(&token)
	if result.Error != nil {
		return nil, result.Error
	}

	return token, nil
}

func (s Service) UpdateToken(ctx *gin.Context, token *models.Token) (*models.Token, error) {

	err := token.Validate()
	if err != nil {
		return nil, err
	}

	tx := s.entityRepository.Db().Begin()
	err = setTransaction(ctx, []*gorm.DB{tx})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	token.IsUsed = &now

	result := tx.Where("uuid=?", token.Uuid).Save(&token)
	if result.Error != nil {
		return nil, result.Error
	}

	return token, nil
}

func (s Service) GetToken(token string) (*models.Token, error) {
	return s.tokenRepository.FindTokenByToken(token)
}

//func (s Service) UpdatePasswordWithToken(ctx context.Context, token *models.Token, exists *models.Token, acc *models.Account, existAcc *models.Account) error {
//
//	_, err := s.tokenRepository.Update(ctx, token, exists)
//	if err != nil {
//		return err
//	}
//
//	_, err = s.accountRepository.Update(ctx, acc, existAcc)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
