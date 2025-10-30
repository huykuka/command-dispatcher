package users

import (
	"command-dispatcher/internal/core/guards"
	"command-dispatcher/internal/core/pipes"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
	//Similar to inject() in Nestjs
	route := r.Group("/users", guards.JWTAuthGuard())

	userService := new(UserService)
	///Register routes
	route.GET("", pipes.Query[GetUserQuery], userService.getAll)
	route.GET("/:id", userService.getByID)
	route.PATCH("/:id", pipes.Body[UpdateUserDTO], userService.update)
	route.POST("", pipes.Body[CreateUserDTO], userService.create)
}
