package server

import (
	"github.com/gin-gonic/gin"
	"password-lock/controller"
	"password-lock/middleware"
)

func initializeRoutes(r *gin.Engine, ctrl *controller.Controller, m *middleware.AuthMiddleware) {
	// user routes
	r.POST("/register", ctrl.RegisterUser)
	r.POST("/verify", ctrl.VerifyAccount)
	r.POST("/login", ctrl.Login)
	r.POST("/forgot-password", ctrl.ForgotPassword)
	r.GET("/security-questions", ctrl.GetSecurityQuestionsByToken)

	r.POST("/complete-registration", ctrl.CompleteRegistration)
	r.POST("/logout", ctrl.Logout)

	// entity routes
	r.POST("/entity", ctrl.CreateEntity)
	r.POST("/entity/update", ctrl.UpdateEntity)
	r.POST("/entity/delete/:entity_uuid", ctrl.DeleteEntity)

	r.GET("/entity/:entity_uuid", ctrl.GetEntity)
	r.GET("/entity/list", ctrl.ListEntities)

	r.POST("/icon/:entity_uuid", ctrl.DownloadEntityIcon)

}
