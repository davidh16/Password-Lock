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

func NewServer(r *gin.Engine, ctrl *controller.Controller, redis *redis.Client) Server {
	middleware := mw.InitializeMiddleware(redis)
	initializeRoutes(r, ctrl, middleware)

	return Server{
		controller: ctrl,
		router:     r,
	}
}

func (s Server) Run(port string) {
	err := s.router.Run(port)
	if err != nil {
		return
	}
}
