package controller

import (
	"github.com/gin-gonic/gin"
	"password-lock/service"
)

type Controller struct {
	service service.Service
}

func NewController(service service.Service) *Controller {
	return &Controller{
		service: service,
	}
}

type Response struct {
	Status  int
	Message string
	Error   string
}

func SendResponse(c *gin.Context, response Response) {
	if len(response.Message) > 0 {
		c.JSON(response.Status, map[string]interface{}{"message": response.Message})
	} else if len(response.Error) > 0 {
		c.JSON(response.Status, map[string]interface{}{"error": response.Error})
	}
}
