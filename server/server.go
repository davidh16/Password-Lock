package server

import (
	"github.com/gin-gonic/gin"
	"password-lock/controller"
)

type Server struct {
	controller *controller.Controller
	router     *gin.Engine
}

func NewServer(r *gin.Engine, ctrl *controller.Controller) Server {
	initializeRoutes(r, ctrl)

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
