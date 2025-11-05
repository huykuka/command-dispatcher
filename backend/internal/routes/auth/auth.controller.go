package auth

import (
	"command-dispatcher/internal/core/guards"
	"command-dispatcher/internal/core/middlewares"
	"command-dispatcher/internal/core/pipes"
	"command-dispatcher/internal/models"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
	route := r.Group("/auth")
	authService := new(AuthService)
	route.POST("/login", middlewares.PublicApiMiddleware(), pipes.Body[models.LoginDTO], authService.Login)
	route.GET("/validate", guards.JWTAuthGuard(), authService.Validate)
}
