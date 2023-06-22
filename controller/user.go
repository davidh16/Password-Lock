package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"password-lock/models"
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

	// TODO generirati secret key

	// TODO poslati secret key na mail

	SendResponse(ctx, Response{
		Status:  http.StatusCreated,
		Message: "user registered successfully",
	})
	return

}
