package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type Middleware struct {
	redis *redis.Client
}

func InitializeMiddleware(redis *redis.Client) *Middleware {
	return &Middleware{
		redis: redis,
	}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionUuid, err := ctx.Cookie("session")
		if err != nil {
			ctx.AbortWithStatus(http.StatusProxyAuthRequired)
			return
		}

		loggedInUser, err := m.redis.Get(context.Background(), sessionUuid).Result()
		if err != nil {
			ctx.AbortWithStatus(http.StatusProxyAuthRequired)
			return
		}

		if loggedInUser == "" {
			ctx.AbortWithStatus(http.StatusProxyAuthRequired)
			return
		} else {
			ctx.Next()
			return
		}
	}
}
