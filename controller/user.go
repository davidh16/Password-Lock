package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"password-lock/models"
	"password-lock/utils"
	"time"
)

func (c Controller) RegisterUser(ctx *gin.Context) {

	var registerRequest struct {
		EmailAddress string `json:"email_address"`
	}

	// decoding json message to user model
	err := json.NewDecoder(ctx.Request.Body).Decode(&registerRequest)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	err, exists := c.service.IfEmailAddressExists(registerRequest.EmailAddress)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	if exists {
		c.SendResponse(ctx, Response{
			Status: http.StatusBadRequest,
			Error:  errors.New("email address already exists").Error(),
		})
		return
	}

	user := &models.User{
		EmailAddress: registerRequest.EmailAddress,
	}

	user, err = c.service.RegisterUser(ctx, user)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	//todo generate verification token
	verificationToken, err := c.service.CreateToken(ctx, user.Uuid, "verification")
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	err = c.service.SendVerificationLinkEmail(user.EmailAddress, verificationToken.Token)
	if err != nil {
		log.Println("failed to send an email")
		return
	}

	c.SendResponse(ctx, Response{
		Status:  http.StatusCreated,
		Message: "user registered successfully",
	})
	return

}

func (c Controller) VerifyAccount(ctx *gin.Context) {

	var verifyAccountRequest struct {
		VerificationToken string `json:"verification_token"`
	}

	// decoding json message to user model
	err := json.NewDecoder(ctx.Request.Body).Decode(&verifyAccountRequest)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	token, err := c.service.GetToken(verifyAccountRequest.VerificationToken)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	if token.IsUsed != nil || token.ExpireAt.Before(time.Now()) {
		c.SendResponse(ctx, Response{
			Status: http.StatusBadRequest,
			Error:  errors.New("invalid token").Error(),
		})
		return
	}

	password, err := utils.GenerateRandomStringURLSafe()
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	user, err := c.service.VerifyUser(ctx, token.UserUuid, password)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	_, err = c.service.UpdateToken(ctx, token)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	err = c.service.SendAccountVerifiedEmail(user.EmailAddress, password)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	c.SendResponse(ctx, Response{
		Status:  http.StatusOK,
		Message: "user successfully verified",
	})
	return
}

func (c Controller) Login(ctx *gin.Context) {

	var credentials models.User

	// decoding json message to user model
	err := json.NewDecoder(ctx.Request.Body).Decode(&credentials)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		})
		return
	}

	user, err := c.service.Authenticate(credentials)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		})
		return
	}

	sessionKey, err := c.service.GenerateAndSaveSessionKey(user)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	ctx.SetCookie("session", sessionKey, 600, "/", "", true, true)

	c.SendResponse(ctx, Response{
		Status:  http.StatusOK,
		Message: "successfully logged in",
	})
	return

}

func (c Controller) Logout(ctx *gin.Context) {
	sessionId, err := ctx.Cookie("session")
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	c.service.TerminateSession(sessionId)

	ctx.SetCookie("session", "", -1, "/", "", true, true)

	c.SendResponse(ctx, Response{
		Status:  http.StatusOK,
		Message: "successfully logged out",
	})
	return
}
