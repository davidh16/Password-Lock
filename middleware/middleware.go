package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"password-lock/db"
	"password-lock/models"
)

type Middleware struct {
	db    *gorm.DB
	redis *redis.Client
}

func InitializeMiddleware(db *gorm.DB, redis *redis.Client) *Middleware {
	return &Middleware{
		db:    db,
		redis: redis,
	}
}

func (m *Middleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionUuid, err := ctx.Cookie("session")
		if err != nil {
			fmt.Println(err)
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

func (m *Middleware) User() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		me, ok := ctx.Get("me")
		if ok {
			result := m.db.Table(db.USERS_TABLE).Where("uuid=?", me).First(&user)
			if result.Error != nil {
				ctx.AbortWithStatus(http.StatusNotFound)
				return
			}
		} else {
			ctx.AbortWithStatus(http.StatusProxyAuthRequired)
			return
		}

		if !user.Completed {
			ctx.AbortWithStatus(http.StatusNotAcceptable)
			return
		}
	}
}

func (m *Middleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5174")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	}
}
