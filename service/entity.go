package service

import "crypto/aes"

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
