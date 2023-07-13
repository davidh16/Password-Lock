package server

import (
	"github.com/gin-gonic/gin"
	"password-lock/controller"
)

func initializeRoutes(r *gin.Engine, ctrl *controller.Controller) {
	r.POST("/register", ctrl.RegisterUser)
}
