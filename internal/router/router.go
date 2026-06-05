package router

import (
	"github.com/asikeida/xiaomajizhang/internal/handler"
	"github.com/asikeida/xiaomajizhang/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth *handler.AuthHandler
	User *handler.UserHandler
}

func New(jwtSecret string, handlers Handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", handlers.Auth.Register)
			authGroup.POST("/login", handlers.Auth.Login)
			authGroup.POST("/refresh", handlers.Auth.Refresh)
			authGroup.POST("/logout", handlers.Auth.Logout)
		}

		v1.GET("/me", middleware.Auth(jwtSecret), handlers.User.Me)
	}

	return r
}
