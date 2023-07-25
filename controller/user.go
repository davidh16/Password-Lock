package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"password-lock/models"
	"password-lock/utils"
)

func (c Controller) RegisterUser(ctx *gin.Context) {
	var user models.User

	// decoding json message to user model
	err := json.NewDecoder(ctx.Request.Body).Decode(&user)
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	_, err = c.service.RegisterUser(user)
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	secretKey, err := utils.GenerateRandomStringURLSafe()
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	err = utils.SendEmail(user.EmailAddress, secretKey)
	if err != nil {
		SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	SendResponse(ctx, Response{
		Status:  http.StatusCreated,
		Message: "user registered successfully",
	})
	return

}
