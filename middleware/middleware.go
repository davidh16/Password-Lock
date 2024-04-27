package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

//type Middleware struct {
//	CORSMiddleware CORSMiddleware
//	AuthMiddleware AuthMiddleware
//}

type CORSMiddleware struct{}

type AuthMiddleware struct {
	redis *redis.Client
}

func InitializeAuthMiddleware(redis *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		redis: redis,
	}
}

func InitializeCORSMiddleware(redis *redis.Client) *CORSMiddleware {
	return &CORSMiddleware{}
}

func (m *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
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
			ctx.Set("me", loggedInUser)
			ctx.Next()
			return
		}
	}
}

func (m *CORSMiddleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	}
}
