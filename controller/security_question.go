package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c Controller) ListSecurityQuestions(ctx *gin.Context) {
	securityQuestions, err := c.service.FindAllSecurityQuestions()
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, securityQuestions)
}
