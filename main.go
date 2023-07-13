package main

import (
	"password-lock/controller"
	"password-lock/db"
	"password-lock/repository"
	"password-lock/server"
	"password-lock/service"

	"github.com/gin-gonic/gin"
)

func main() {

	pgInstance := db.ConnectToDatabase()

	repo := repository.NewRepository(pgInstance)

	svc := service.NewService(repo)

	ctrl := controller.NewController(svc)
	r := gin.Default()

	srv := server.NewServer(r, ctrl)

	// Listen and Server in 0.0.0.0:8080
	srv.Run(":8080")
}
