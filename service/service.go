package service

import (
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/option"
	"password-lock/config"
	"password-lock/repository"
)

type Service struct {
	redis            *redis.Client
	Cfg              *config.Config
	userRepository   repository.UserRepository
	entityRepository repository.EntityRepository
	tokenRepository  repository.TokenRepository
	firebaseApp      *firebase.App
}

func NewService(redis *redis.Client,
	config *config.Config,
	userRepo repository.UserRepository,
	entityRepo repository.EntityRepository,
	tokenRepo repository.TokenRepository,
) Service {

	opt := option.WithCredentialsJSON([]byte(config.FirebaseCredentialsJSON))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Println("Could not initialize Firebase client : ", err.Error())
	}

	return Service{
		redis:            redis,
		Cfg:              config,
		userRepository:   userRepo,
		entityRepository: entityRepo,
		tokenRepository:  tokenRepo,
		firebaseApp:      app,
	}
}
