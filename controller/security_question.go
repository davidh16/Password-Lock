package controller

import (
	"encoding/json"
	"fmt"
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

	data := map[string]interface{}{
		"security_questions": securityQuestions,
	}

	response, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	ctx.JSON(http.StatusOK, string(response))
}
