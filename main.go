package main

import (
	"password-lock/config"
	"password-lock/controller"
	"password-lock/db"
	mw "password-lock/middleware"
	"password-lock/repository"
	"password-lock/server"
	"password-lock/service"
)

func main() {

	cfg := config.GetConfig()

	redis := db.ConnectToRedis(cfg)

	pgInstance := db.ConnectToDatabase(cfg)

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

	middleware := mw.InitializeMiddleware(pgInstance, redis)

	srv := server.NewServer(ctrl, middleware, cfg)

	// Listen and Server in 0.0.0.0:8080
	srv.Run(":8080")
}
