package main

import (
	"password-lock/config"
	"password-lock/controller"
	"password-lock/db"
	"password-lock/repository"
	"password-lock/server"
	"password-lock/service"
)

func main() {

	cfg := config.GetConfig()

	redis := db.ConnectToRedis()

	pgInstance := db.ConnectToDatabase()

	userRepo := repository.NewUserRepository(pgInstance)
	entityRepo := repository.NewEntityRepository(pgInstance)
	tokenRepo := repository.NewTokenRepository(pgInstance)

	svc := service.NewService(
		redis,
		cfg,
		userRepo,
		entityRepo,
		tokenRepo,
	)

	ctrl := controller.NewController(svc)

	srv := server.NewServer(ctrl, redis)

	// Listen and Server in 0.0.0.0:8080
	srv.Run(":8080")
}
