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

	// in the future, this will fetch icon paths from some kind of cloud storage
	if entity.Type != 5 {
		iconPath := c.service.GetEntityIconPath(entity.Type)
		entity.IconPath = &iconPath
	}

	entity.Password = c.service.EncryptPassword(entity.SecretKey, entity.Password)

	_, err = c.service.CreateEntity(entity)
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

}
