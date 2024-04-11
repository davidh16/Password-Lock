package server

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"password-lock/controller"
	mw "password-lock/middleware"
)

type Server struct {
	controller *controller.Controller
	router     *gin.Engine
}

func NewServer(ctrl *controller.Controller, redis *redis.Client) Server {

	middleware := mw.InitializeMiddleware(redis)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	initializeRoutes(router, ctrl, middleware)

	return Server{
		controller: ctrl,
		router:     router,
	}
}

func (s Server) Run(port string) {
	err := s.router.Run(port)
	if err != nil {
		return
	}
}
