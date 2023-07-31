package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"password-lock/models"
)

func (c Controller) CreateEntity(ctx *gin.Context) {
	var entity models.Entity

	// decoding json message to user model
	err := json.NewDecoder(ctx.Request.Body).Decode(&entity)
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	if entity.Password == "" {
		SendResponse(ctx, Response{
			Status: http.StatusBadRequest,
			Error:  "password not provided for given entity",
		})
		return
	}

	if entity.SecretKey == "" {
		SendResponse(ctx, Response{
			Status: http.StatusBadRequest,
			Error:  "secret key not provided",
		})
		return
	}

	//TODO pozvati servise da na temelju tipa entiteta dohvati ikonicu sa diska ili bucketa, odnosno path

	encryptedPassword := c.service.EncryptPassword(entity.SecretKey, entity.Password)

}
