package command

import (
	//"command-dispatcher/internal/config/_queue"
	//"command-dispatcher/internal/worker"
	//"fmt"

	"command-dispatcher/internal/models"
	"command-dispatcher/internal/worker"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CommandService struct{}

func NewCommandService() *CommandService {
	return &CommandService{}
}

func (s *CommandService) add(c *gin.Context) {
	dto := c.MustGet("Body").(models.CommandCreateDTO)
	err := worker.EnqueueCommandExecutionTask(dto)
	if err != nil {
		fmt.Println("Could not enqueue task: ", err)
		return
	}
}

func (s *CommandService) remove(c *gin.Context) {
	// Implement command removal logic here

}

func (s *CommandService) list() {
	// Implement command listing logic here

}

func (s *CommandService) get(commandID string) {
	// Implement command retrieval logic here

}

func (s *CommandService) update(commandID string, newCommand string) {
	// Implement command update logic here

}

func (s *CommandService) execute(commandID string) {
	// Implement command execution logic here

}

func (s *CommandService) schedule(commandID string, scheduleTime string) {
	// Implement command scheduling logic here

}

func (s *CommandService) cancel(commandID string) {
}
