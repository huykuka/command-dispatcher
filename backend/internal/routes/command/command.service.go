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
// @Summary Create a new command configuration
// @Description Create a new command configuration with the provided details
// @Tags commands
// @Accept json
// @Produce json
// @Param command body models.CommandConfigCreateDTO true "Command Config"
// @Success 201 {object} db.CommandConfig
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /command [post]
func (s *CommandService) create(c *gin.Context) {

	dto := c.MustGet("Body").(models.CommandConfigCreateDTO)

	// Convert DTO to database model
	commandConfig := dto.ToEntity()

	if err := s.repo.Create(commandConfig); err != nil {
		utils.HandleHTTPError(c, "Create command config failed", "Create command config failed")
		return
	}

	c.Status(201)
	c.Set("response", commandConfig)
}

// getAll retrieves all command configurations.
// @Summary Get all command configurations
// @Description Retrieve a list of all command configurations
// @Tags commands
// @Produce json
// @Success 200 {array} db.CommandConfig
// @Failure 500 {object} map[string]interface{}
// @Router /command [get]
func (s *CommandService) getAll(c *gin.Context) {
	commands, err := s.repo.FindAll()
	if err != nil {
		utils.HandleHTTPError(c, "Create command config failed", "Fetch command configs failed")
		return
	}
	c.Status(200)
	c.Set("response", commands)
}

// getByID retrieves a single command configuration by its ID.
// @Summary Get command configuration by ID
// @Description Retrieve a specific command configuration by its ID
// @Tags commands
// @Produce json
// @Param id path string true "Command Config ID"
// @Success 200 {object} db.CommandConfig
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /command/{id} [get]
func (s *CommandService) getByID(c *gin.Context) {
	id := c.Param("id")
	command, err := s.repo.FindByID(id)
	if err != nil {
		utils.HandleHTTPError(c, "Create command config failed", "Fetch command config failed")
		return
	}
	c.Set("response", command)
}

// update updates an existing command configuration.
// @Summary Update command configuration
// @Description Update an existing command configuration with partial data
// @Tags commands
// @Accept json
// @Produce json
// @Param id path string true "Command Config ID"
// @Param command body models.CommandConfigUpdateDTO true "Command Config Update"
// @Success 200 {object} db.CommandConfig
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /command/{id} [patch]
func (s *CommandService) update(c *gin.Context) {
	id := c.Param("id")
	dto := c.MustGet("Body").(models.CommandConfigUpdateDTO)
	command, err := s.repo.FindByID(id)
	if err != nil {
		utils.HandleHTTPError(c, "Create command config failed", "Fetch command config failed")
		return
	}

	// Apply DTO updates to entity
	dto.ApplyTo(command)

	if err := s.repo.Update(command); err != nil {
		utils.HandleHTTPError(c, "Create command config failed", "Update command config failed")
		return
	}
	c.Set("response", command)
}

// delete deletes a command configuration.
// @Summary Delete command configuration
// @Description Delete a command configuration by ID
// @Tags commands
// @Param id path string true "Command Config ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /command/{id} [delete]
func (s *CommandService) delete(c *gin.Context) {
	id := c.Param("id")
	if err := s.repo.Delete(id); err != nil {
		utils.HandleHTTPError(c, "Create command config failed", "Delete command config failed")
		return
	}
	c.Status(204)
}
