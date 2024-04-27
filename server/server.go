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

	authMw := mw.InitializeAuthMiddleware(redis)
	_ = mw.InitializeCORSMiddleware(redis)

	router := gin.Default()

	initializeRoutes(router, ctrl, authMw)

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
