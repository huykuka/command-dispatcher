package command

import (
	"command-dispatcher/internal/core/pipes"
	"command-dispatcher/internal/models"

	"github.com/gin-gonic/gin"
)

// Register sets up the command routes within the provided Gin router group.
func Register(r *gin.RouterGroup) {
	route := r.Group("/command")

	commandService := NewCommandService()
	
	route.POST("", pipes.Body[models.CommandConfigCreateDTO], commandService.create)
	route.GET("", commandService.getAll)
	route.GET("/:id", commandService.getByID)
	route.PATCH("/:id", pipes.Body[models.CommandConfigUpdateDTO], commandService.update)
	route.DELETE("/:id", commandService.delete)
}
