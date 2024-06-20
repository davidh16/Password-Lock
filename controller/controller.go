package controller

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"password-lock/service"
)

type Controller struct {
	service service.Service
}

func NewController(service service.Service) *Controller {
	return &Controller{
		service: service,
	}
}

type Response struct {
	Status  int
	Message string
	Error   string
}

func (c Controller) SendResponse(ctx *gin.Context, response Response) {

	xxx, _ := ctx.Get("tx")
	transactions, _ := xxx.([]*gorm.DB)

	if len(transactions) > 0 {
		if response.Error != "" {
			for _, tx := range transactions {
				tx.Rollback()
			}
		} else {
			for _, tx := range transactions {
				tx.Commit()
			}
		}
	}

	if len(response.Message) > 0 {
		ctx.JSON(response.Status, map[string]string{"message": response.Message})
		return
	} else if len(response.Error) > 0 {
		ctx.JSON(response.Status, map[string]string{"error": response.Error})
		return
	}

	ctx.Writer.WriteHeader(response.Status)
}

func (c Controller) encryptResponse(text string) (string, error) {

	iv, err := hex.DecodeString(c.service.Cfg.ResponseSecretVector)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(c.service.Cfg.ResponseSecretKey))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return service.Encode(cipherText), nil
}

func ignoreNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return err
	}
}
