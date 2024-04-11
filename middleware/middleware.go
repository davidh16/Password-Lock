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

func (m *Middleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
