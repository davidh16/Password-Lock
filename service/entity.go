package service

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"password-lock/models"
	"strings"
)

func (s Service) EncryptPassword(secretKey string, password string) string {

	c, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())

	encryptedPassowrd := gcm.Seal(nonce, nonce, []byte(password), nil)

	return base64.StdEncoding.EncodeToString(encryptedPassowrd)
}

func (s Service) GetEntityIconPath(entityType int) string {

	iconType := models.TypeMap[entityType]
	if iconType != "custom" {
		return strings.Join([]string{"default", iconType}, "/")
	}
	return ""
}

func (s Service) CreateEntity(ctx *gin.Context, entity models.Entity) (*models.Entity, error) {
	tx := s.entityRepository.Db().Begin()
	err := setTransaction(ctx, []*gorm.DB{tx})
	if err != nil {
		return nil, err
	}

	result := tx.Create(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func (s Service) UpdateEntity(ctx *gin.Context, updatedEntity *models.Entity) error {
	tx := s.entityRepository.Db().Begin()
	err := setTransaction(ctx, []*gorm.DB{tx})
	if err != nil {
		return err
	}

	var entity models.Entity
	result := tx.Where("uuid=?", updatedEntity.Uuid).First(&entity)
	if result.Error != nil {
		return result.Error
	}

	entity.Merge(updatedEntity)

	result = tx.Where("uuid=?", updatedEntity.Uuid).Save(&entity)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s Service) DeleteEntity(ctx *gin.Context, entityUuid string) error {
	tx := s.entityRepository.Db().Begin()
	err := setTransaction(ctx, []*gorm.DB{tx})
	if err != nil {
		return err
	}

	result := tx.Where("uuid=?", entityUuid).Delete(models.Entity{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s Service) GetEntityByUuid(ctx *gin.Context, entityUuid string, secretKey string) (*models.Entity, error) {

	me := ctx.Value("me").(string)

	var entity models.Entity
	result := s.entityRepository.Db().Where("uuid=? AND user_uuid=?", entityUuid, me).First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}

	decryptedPassword, err := decryptEntityPassword(entity.Password, secretKey)
	if err != nil {
		return nil, err
	}

	entity.Password = decryptedPassword

	return &entity, nil
}

func (s Service) ListEntities(ctx *gin.Context) ([]models.Entity, error) {
	var entities []models.Entity
	me := ctx.Value("me").(string)
	result := s.entityRepository.Db().Where("user_uuid=?", me).Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}

	return hideEntityPasswords(entities), nil
}

func hideEntityPasswords(entities []models.Entity) []models.Entity {
	var entitiyListWithHiddenPasswords []models.Entity
	for _, entity := range entities {
		entitiyListWithHiddenPasswords = append(entitiyListWithHiddenPasswords, *entity.HidePassword())
	}
	return entitiyListWithHiddenPasswords
}

func decryptEntityPassword(password string, secretKey string) (string, error) {

	ciphertext, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return "", err
	}

	key := []byte(secretKey)

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	decryptedPassword, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(decryptedPassword), nil

}
