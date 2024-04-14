package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func (c Controller) SendResponse(ctx *gin.Context, response Response) {

	xxx, _ := ctx.Get("tx")
	transactions, _ := xxx.([]*gorm.DB)

	fmt.Println(len(transactions))

	if len(transactions) > 0 {
		if response.Error != "" {
			for _, tx := range transactions {
				tx.Rollback()
			}
		} else {
			for _, tx := range transactions {
				tx.Commit()
			}
		}
	}

	if len(response.Message) > 0 {
		ctx.JSON(response.Status, map[string]interface{}{"message": response.Message})
	} else if len(response.Error) > 0 {
		ctx.JSON(response.Status, map[string]interface{}{"error": response.Error})
	}
}
