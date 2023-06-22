package service

import "github.com/gin-gonic/gin"

func (s Service) RegisterUser(ctx *gin.Context) {
	s.repository.Testiram()
}
