package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"password-lock/config"
	"password-lock/repository"
)

type Service struct {
	redis            *redis.Client
	Cfg              *config.Config
	userRepository   repository.UserRepository
	entityRepository repository.EntityRepository
	tokenRepository  repository.TokenRepository
}

func NewService(redis *redis.Client,
	config *config.Config,
	userRepo repository.UserRepository,
	entityRepo repository.EntityRepository,
	tokenRepo repository.TokenRepository,
) Service {
	return Service{
		redis:            redis,
		Cfg:              config,
		userRepository:   userRepo,
		entityRepository: entityRepo,
		tokenRepository:  tokenRepo,
	}
}

func setTransaction(ctx *gin.Context, transactionsToCommit []*gorm.DB) error {
	tx, exists := ctx.Get("tx")
	if exists {
		transactions, ok := tx.([]*gorm.DB)
		if ok {
			transactions = append(transactions, transactionsToCommit...)
			ctx.Set("tx", transactions)
			return nil
		} else {
			return errors.New("something went wrong")
		}
	} else {
		ctx.Set("tx", transactionsToCommit)
		return nil
	}
}
