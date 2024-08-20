package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"password-lock/config"
	"password-lock/db"
	"password-lock/models"
)

type Middleware struct {
	db    *gorm.DB
	redis *redis.Client
	cfg   *config.Config
}

func InitializeMiddleware(db *gorm.DB, redis *redis.Client, cfg *config.Config) *Middleware {
	return &Middleware{
		db:    db,
		redis: redis,
		cfg:   cfg,
	}
}

func (m *Middleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionUuid, err := ctx.Cookie("session")
		if err != nil {
			ctx.JSON(http.StatusProxyAuthRequired, "session expired")
			//ctx.AbortWithError(http.StatusProxyAuthRequired, errors.New("session expired"))
			return
		}

		loggedInUser, err := m.redis.Get(context.Background(), sessionUuid).Result()
		if err != nil {
			ctx.AbortWithError(http.StatusProxyAuthRequired, errors.New("session expired"))
			return
		}

		if loggedInUser == "" {
			ctx.AbortWithError(http.StatusProxyAuthRequired, errors.New("session expired"))
			return
		} else {
			ctx.Set("me", loggedInUser)
			ctx.Next()
			return
		}
	}
}

func (m *Middleware) Session() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		sessionUuid, err := ctx.Cookie("session")
		if err != nil {
			data := map[string]interface{}{
				"is_authenticated": false,
			}

			response, err := json.Marshal(data)
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			ctx.JSON(403, response)
		}

		var isAuthenticated bool
		var status int

		loggedInUser, err := m.redis.Get(context.Background(), sessionUuid).Result()
		if loggedInUser == "" || err != nil {
			isAuthenticated = false
			status = http.StatusForbidden
		} else {
			isAuthenticated = true
			status = http.StatusOK
		}

		data := map[string]interface{}{
			"is_authenticated": isAuthenticated,
		}
		response, err := json.Marshal(data)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(status, response)
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
		c.Writer.Header().Set("Access-Control-Allow-Origin", m.cfg.FrontendBaseUrl)
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
