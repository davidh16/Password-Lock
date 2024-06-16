package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"password-lock/models"
	"password-lock/utils"
	"password-lock/validations"
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

	verificationToken, err := c.service.CreateToken(ctx, user.Uuid, models.TOKEN_TYPE_VERIFICATION)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	err = c.service.SendVerificationLinkEmail(user.EmailAddress, verificationToken.Token)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusBadRequest,
			Error:  errors.New("email address is not valid").Error(),
		})
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

	tokenString := ctx.Query("token")

	token, err := c.service.GetToken(tokenString)
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

	err = c.service.SendPasswordEmail(user.EmailAddress, password)
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

func (c Controller) ResendVerificationEmail(ctx *gin.Context) {
	var resendVerificationEmailRequest struct {
		EmailAddress string `json:"email_address"`
	}

	err := json.NewDecoder(ctx.Request.Body).Decode(&resendVerificationEmailRequest)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
		})
		return
	}

	user, err := c.service.GetUserByEmailAddress(resendVerificationEmailRequest.EmailAddress)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.SendResponse(ctx, Response{
				Status: http.StatusInternalServerError,
			})
			return
		}
	}

	user, err = c.service.GetUnverifiedUserByEmailAddress(resendVerificationEmailRequest.EmailAddress)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
		})
		return
	}

	token, err := c.service.CreateToken(ctx, user.Uuid, models.TOKEN_TYPE_VERIFICATION)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
		})
		return
	}

	err = c.service.SendVerificationLinkEmail(user.EmailAddress, token.Token)
	if err != nil {
		log.Println("failed to send an email")
		return
	}

	c.SendResponse(ctx, Response{
		Status: http.StatusOK,
	})
	return
}

func (c Controller) CompleteRegistration(ctx *gin.Context) {
	me, err := c.service.Me(ctx)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	if me.Completed {
		c.SendResponse(ctx, Response{
			Status: http.StatusBadRequest,
			Error:  errors.New("registration already completed").Error(),
		})
		return
	}

	var userPersonalQuestions []*models.UserPersonalQuestion
	err = json.NewDecoder(ctx.Request.Body).Decode(&userPersonalQuestions)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	for i, _ := range userPersonalQuestions {
		userPersonalQuestions[i].UserUuid = me.Uuid
	}

	if !validations.IsCompleteRegistrationRequestValid(userPersonalQuestions) {
		c.SendResponse(ctx, Response{
			Status: http.StatusBadRequest,
			Error:  errors.New("all fields must be populated").Error(),
		})
		return
	}

	me.Completed = true

	_, err = c.service.CompleteRegistration(ctx, me, userPersonalQuestions)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	c.SendResponse(ctx, Response{
		Status:  http.StatusOK,
		Message: "registration completed successfully",
	})
	return
}

func (c Controller) Login(ctx *gin.Context) {

	var credentials struct {
		EmailAddress string `json:"email_address"`
		Password     string `json:"password"`
		RememberMe   bool   `json:"remember_me"`
	}

	// decoding json message to user model
	err := json.NewDecoder(ctx.Request.Body).Decode(&credentials)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		})
		return
	}

	user, err := c.service.Authenticate(credentials.EmailAddress, credentials.Password)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		})
		return
	}

	sessionLifeTime := int(time.Hour)
	if credentials.RememberMe {
		sessionLifeTime = int(time.Hour * 24)
	}

	sessionKey, err := c.service.GenerateAndSaveSessionKey(user.Uuid, time.Duration(sessionLifeTime))
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	ctx.SetCookie("session", sessionKey, sessionLifeTime, "/", "", true, false)

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

func (c Controller) ForgotPassword(ctx *gin.Context) {
	var forgotPasswordRequest struct {
		EmailAddress string `json:"email_address"`
	}

	// decoding json message to user model
	err := json.NewDecoder(ctx.Request.Body).Decode(&forgotPasswordRequest)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		})
		return
	}

	// if user is not found, we can not return that information because of security issues
	user, err := c.service.GetUserByEmailAddress(forgotPasswordRequest.EmailAddress)
	if err != nil {
		return
	}

	token, err := c.service.CreateToken(ctx, user.Uuid, models.TOKEN_TYPE_FORGOT_PASSWORD)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	err = c.service.SendPasswordResetLinkEmail(user.EmailAddress, token.Token)
	if err != nil {
		log.Println("failed to send an email")
		return
	}

	c.SendResponse(ctx, Response{
		Status: http.StatusOK,
	})
	return
}

func (c Controller) GetUserPersonalQuestionsByToken(ctx *gin.Context) {
	t := ctx.Query("token")

	token, err := c.service.GetToken(t)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	if token == nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusForbidden,
		})
		return
	}

	if token.IsUsed != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusForbidden,
		})
		return
	}

	if token.TokenType != models.TOKEN_TYPE_FORGOT_PASSWORD {
		c.SendResponse(ctx, Response{
			Status: http.StatusForbidden,
		})
		return
	}

	ctx.JSON(http.StatusOK, token.User.PersonalQuestions)
	return
}

func (c Controller) CheckPersonalQuestionsAnswers(ctx *gin.Context) {
	t := ctx.Query("token")

	token, err := c.service.GetToken(t)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	if token == nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusForbidden,
		})
		return
	}

	if token.IsUsed != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusForbidden,
		})
		return
	}

	if token.TokenType != models.TOKEN_TYPE_FORGOT_PASSWORD {
		c.SendResponse(ctx, Response{
			Status: http.StatusForbidden,
		})
		return
	}

	var personalQuestionsFromRequest []models.UserPersonalQuestion
	err = json.NewDecoder(ctx.Request.Body).Decode(&personalQuestionsFromRequest)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	err = validations.ValidatePersonalQuestionsAnswers(token.User.PersonalQuestions, personalQuestionsFromRequest)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusForbidden,
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

	passswordResetToken, err := c.service.CreateToken(ctx, token.UserUuid, models.TOKEN_TYPE_PASSWORD_RESET)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	ctx.JSON(200, struct {
		PasswordResetToken string `json:"password_reset_token"`
	}{PasswordResetToken: passswordResetToken.Token})
	return
}

func (c Controller) ResetPassword(ctx *gin.Context) {
	t := ctx.Query("token")

	token, err := c.service.GetToken(t)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	if token == nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusForbidden,
		})
		return
	}

	if token.IsUsed != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusForbidden,
		})
		return
	}

	if token.TokenType != models.TOKEN_TYPE_PASSWORD_RESET {
		c.SendResponse(ctx, Response{
			Status: http.StatusForbidden,
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

	newPassword, err := utils.GenerateRandomStringURLSafe()
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	err = c.service.UpdatePassword(ctx, &token.User, newPassword)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	err = c.service.SendNewPasswordEmail(token.User.EmailAddress, newPassword)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	c.SendResponse(ctx, Response{
		Status: http.StatusOK,
	})
	return
}

func (c Controller) Me(ctx *gin.Context) {
	me, err := c.service.Me(ctx)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, me)
}
