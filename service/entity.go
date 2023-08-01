package service

import (
	"crypto/aes"
	"password-lock/models"
)

func (s Service) EncryptPassword(secretKey string, password string) string {

	cipher, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		panic(err)
	}

	// Make a buffer the same length as plaintext
	encryptedPassword := make([]byte, len(password))
	cipher.Encrypt(encryptedPassword, []byte(password))

	return string(encryptedPassword)
}

func (s Service) GetEntityIconPath(entityType int) string {
	switch entityType {
	case 0:
		return "/Users/davidhorvat/GolandProjects/Password-lock/logos/github.png"
	case 1:
		return "/Users/davidhorvat/GolandProjects/Password-lock/logos/facebook.png"
	case 2:
		return "/Users/davidhorvat/GolandProjects/Password-lock/logos/gmail.png"
	case 3:
		return "/Users/davidhorvat/GolandProjects/Password-lock/logos/linkedin.png"
	case 4:
		return "/Users/davidhorvat/GolandProjects/Password-lock/logos/instagram.png"
	}
	return ""
}

func (s Service) CreateEntity(entity models.Entity) (*models.Entity, error) {
	result := s.entityRepository.Db().Create(entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}
