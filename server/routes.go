package server

import (
	"github.com/gin-gonic/gin"
	"password-lock/controller"
	"password-lock/middleware"
)

func initializeRoutes(r *gin.Engine, ctrl *controller.Controller, m *middleware.Middleware) {
	// user routes
	r.POST("/register", ctrl.RegisterUser)
	r.POST("/verify", ctrl.VerifyAccount)
	r.POST("/login", ctrl.Login)
	r.POST("/forgot-password", ctrl.ForgotPassword)
	r.POST("/personal-questions", ctrl.GetUserPersonalQuestionsByToken)
	r.POST("/check-personal-questions", ctrl.CheckPersonalQuestionsAnswers)
	r.POST("/reset-password", ctrl.ResetPassword)
	r.POST("/resend-verification-email", ctrl.ResendVerificationEmail)
	r.GET("/list-security-questions", ctrl.ResendVerificationEmail)

	r.POST("/me", m.Auth(), ctrl.Me)
	r.POST("/complete-registration", m.Auth(), ctrl.CompleteRegistration)
	r.POST("/logout", m.Auth(), m.User(), ctrl.Logout)

	// entity routes
	r.POST("/entity", m.Auth(), m.User(), ctrl.CreateEntity)
	r.POST("/entity/update", m.Auth(), m.User(), ctrl.UpdateEntity)
	r.POST("/entity/delete/:entity_uuid", m.Auth(), m.User(), ctrl.DeleteEntity)

	r.GET("/entity/:entity_uuid", m.Auth(), m.User(), ctrl.GetEntity)
	r.GET("/entity/list", m.Auth(), m.User(), ctrl.ListEntities)

	r.POST("/icon/:entity_uuid", m.Auth(), m.User(), ctrl.DownloadEntityIcon)

}
