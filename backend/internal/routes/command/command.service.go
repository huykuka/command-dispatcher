package command

import (
	"command-dispatcher/internal/config/db"
	"command-dispatcher/internal/models"
	"command-dispatcher/internal/utils"

	"github.com/gin-gonic/gin"
)

// CommandService provides command-related business logic.
type CommandService struct {
	repo *CommandRepository
}

// NewCommandService creates a new CommandService instance.
func NewCommandService() *CommandService {
	database := db.GetDB()
	return &CommandService{repo: NewCommandRepository(database)}
}

// create handles creating a new command configuration.
func (s *CommandService) create(c *gin.Context) {
	
	dto := c.MustGet("Body").(models.CommandConfigCreateDTO)
	command := &db.CommandConfig{
		Name:                   dto.Name,
		Description:            dto.Description,
		IsAcknowledgeRequired:  dto.IsAcknowledgeRequired,
		CommandType:            dto.CommandType,
		PayloadSchema:          dto.PayloadSchema,
		AknowlegmentTimeout:    dto.AknowlegmentTimeout,
		CompletionTimeout:      dto.CompletionTimeout,
	}

	if err := s.repo.Create(command); err != nil {
		utils.HandleHTTPError(c, err.Error(), "Create command config failed")
		return
	}
	c.Status(201)
	utils.SetResponse(c, map[string]any{
		"id":          command.ID,
		"name":        command.Name,
		"description": command.Description,
	})
}

// getAll retrieves all command configurations.
func (s *CommandService) getAll(c *gin.Context) {
	commands, err := s.repo.FindAll()
	if err != nil {
		utils.HandleHTTPError(c, err.Error(), "Fetch command configs failed")
		return
	}
	c.Status(200)
	c.Set("response", commands)
}

// getByID retrieves a single command configuration by its ID.
func (s *CommandService) getByID(c *gin.Context) {
	id := c.Param("id")
	command, err := s.repo.FindByID(id)
	if err != nil {
		utils.HandleHTTPError(c, err.Error(), "Fetch command config failed")
		return
	}
	c.Set("response", command)
}

// update updates an existing command configuration.
func (s *CommandService) update(c *gin.Context) {
	id := c.Param("id")
	dto := c.MustGet("Body").(models.CommandConfigUpdateDTO)
	command, err := s.repo.FindByID(id)
	if err != nil {
		utils.HandleHTTPError(c, err.Error(), "Fetch command config failed")
		return
	}

	// Update fields from DTO
	if dto.Name != nil {
		command.Name = *dto.Name
	}
	if dto.Description != nil {
		command.Description = *dto.Description
	}
	if dto.CommandType != nil {
		command.CommandType = *dto.CommandType
	}
	if dto.IsAcknowledgeRequired != nil {
		command.IsAcknowledgeRequired = *dto.IsAcknowledgeRequired
	}
	if dto.PayloadSchema != nil {
		command.PayloadSchema = *dto.PayloadSchema
	}
	if dto.AknowlegmentTimeout != nil {
		command.AknowlegmentTimeout = *dto.AknowlegmentTimeout
	}
	if dto.CompletionTimeout != nil {
		command.CompletionTimeout = *dto.CompletionTimeout
	}

	if err := s.repo.Update(command); err != nil {
		utils.HandleHTTPError(c, err.Error(), "Update command config failed")
		return
	}
	c.Set("response", command)
}

// delete deletes a command configuration.
func (s *CommandService) delete(c *gin.Context) {
	id := c.Param("id")
	if err := s.repo.Delete(id); err != nil {
		utils.HandleHTTPError(c, err.Error(), "Delete command config failed")
		return
	}
	c.Status(204)
}
