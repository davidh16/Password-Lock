package routes

import (
	"github.com/gin-gonic/gin"
	"password-lock/controller"
)

func InitializeRoutes(r *gin.Engine, c controller.Controller) {
	r.GET("/nesto", c.RegisterUser)
}
