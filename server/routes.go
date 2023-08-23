package server

import (
	"github.com/gin-gonic/gin"
	"password-lock/controller"
	"password-lock/middleware"
)

func initializeRoutes(r *gin.Engine, ctrl *controller.Controller, m *middleware.Middleware) {
	// user routes
	r.POST("/register", ctrl.RegisterUser)
	r.POST("/login", ctrl.Login)
	r.POST("/logout", ctrl.Logout)

	// entity routes
	r.Use(m.AuthMiddleware()).POST("/entity", ctrl.CreateEntity)
	r.Use(m.AuthMiddleware()).POST("/entity/update", ctrl.UpdateEntity)
	r.Use(m.AuthMiddleware()).POST("/entity/delete/:entity_uuid", ctrl.DeleteEntity)
	r.Use(m.AuthMiddleware()).POST("/entity/:entity_uuid", ctrl.GetEntity)
	r.Use(m.AuthMiddleware()).GET("/entity/list", ctrl.ListEntities)

}
