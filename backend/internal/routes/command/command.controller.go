package command

import (
	"command-dispatcher/internal/core/pipes"
	"command-dispatcher/internal/models"

	// Import the tasks package
	// For logging errors during task enqueuing
	"github.com/gin-gonic/gin"
	// Import asynq for client operations
)

func Register(r *gin.RouterGroup) {
	route := r.Group("/command")

	commandService := NewCommandService()
	route.POST("", pipes.Body[models.CommandCreateDTO], commandService.add) //
}
