package server

import (
	"github.com/gin-gonic/gin"
	"password-lock/controller"
)

func initializeRoutes(r *gin.Engine, ctrl *controller.Controller) {
	// user routes
	r.POST("/register", ctrl.RegisterUser)
	r.POST("/login", ctrl.Login)
	r.POST("/logout", ctrl.Logout)

	// entity routes
	r.POST("/entity", ctrl.CreateEntity)
	r.POST("/entity/delete/:entity_uuid", ctrl.DeleteEntity)
	r.POST("/entity/:entity_uuid", ctrl.GetEntity)
}
