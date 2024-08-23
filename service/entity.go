package service

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"password-lock/models"
	"reflect"
	"strings"
)

func (s Service) GetEntityIconPath(entityType int) string {

	iconType := models.TypeMap[entityType]
	if iconType != "custom" {
		return strings.Join([]string{"default", iconType + ".png"}, "/")
	}
	return ""
}

func (s Service) CreateEntity(ctx *gin.Context, entity models.Entity) (*models.Entity, error) {

	encryptedPassword, err := s.encryptEntity(entity.Password)
	if err != nil {
		return nil, err
	}

	entity.Password = encryptedPassword

	result := s.entityRepository.Db().Create(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func (s Service) UpdateEntity(ctx *gin.Context, updatedEntity *models.Entity) error {

	var entity models.Entity
	result := s.entityRepository.Db().Where("uuid=?", updatedEntity.Uuid).First(&entity)
	if result.Error != nil {
		return result.Error
	}

	if updatedEntity.Type != entity.Type {
		entity.IconPath = s.GetEntityIconPath(updatedEntity.Type)
	}

	entity.Merge(updatedEntity)

	decryptedPassword, err := s.decryptEntity(entity.Password)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(updatedEntity.Password, decryptedPassword) {
		encryptedPassword, err := s.encryptEntity(updatedEntity.Password)
		if err != nil {
			return err
		}

		entity.Password = encryptedPassword
	}

	result = s.entityRepository.Db().Where("uuid=?", updatedEntity.Uuid).Save(&entity)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s Service) DeleteEntity(ctx *gin.Context, entityUuid string) error {
	result := s.entityRepository.Db().Where("uuid=?", entityUuid).Delete(models.Entity{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s Service) GetEntityByUuid(ctx *gin.Context, entityUuid string) (*models.Entity, error) {

	me := ctx.Value("me").(string)

	var entity models.Entity
	result := s.entityRepository.Db().Where("uuid=? AND user_uuid=?", entityUuid, me).First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}

	decryptedPassword, err := s.decryptEntity(entity.Password)
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

	for i, _ := range entities {
		decryptedPassword, err := s.decryptEntity(entities[i].Password)
		if err != nil {
			return nil, err
		}
		entities[i].Password = decryptedPassword
	}

	return entities, nil
}

func (s Service) encryptEntity(text string) (string, error) {

	iv, err := hex.DecodeString(s.Cfg.EntitySecretVector)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(s.Cfg.EntitySecretKey))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func (s Service) decryptEntity(text string) (string, error) {

	iv, err := hex.DecodeString(s.Cfg.EntitySecretVector)
	if err != nil {
		return "", err
	}

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher([]byte(s.Cfg.EntitySecretKey))
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}
