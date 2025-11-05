package users

import (
	"command-dispatcher/internal/core/guards"
	"command-dispatcher/internal/core/pipes"
	"command-dispatcher/internal/models"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
	//Similar to inject() in Nestjs
	route := r.Group("/users", guards.JWTAuthGuard())

	userService := new(UserService)
	///Register routes
	route.GET("", pipes.Query[models.GetUserQuery], userService.getAll)
	route.GET("/:id", userService.getByID)
	route.PATCH("/:id", pipes.Body[models.UpdateUserDTO], userService.update)
	route.POST("", pipes.Body[models.CreateUserDTO], userService.create)
}
