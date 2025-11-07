package command

import (
	"command-dispatcher/internal/config/db"

	"gorm.io/gorm"
)

type CommandRepository struct {
	db *gorm.DB
}

func NewCommandRepository(database *gorm.DB) *CommandRepository {
	return &CommandRepository{db: database}
}

func (r *CommandRepository) Create(command *db.CommandConfig) error {
	return r.db.Create(command).Error
}

func (r *CommandRepository) FindAll() ([]db.CommandConfig, error) {
	var commands []db.CommandConfig
	if err := r.db.Find(&commands).Error; err != nil {
		return nil, err
	}
	return commands, nil
}

func (r *CommandRepository) FindByID(id string) (*db.CommandConfig, error) {
	var command db.CommandConfig
	if err := r.db.First(&command, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &command, nil
}

func (r *CommandRepository) Update(command *db.CommandConfig) error {
	return r.db.Save(command).Error
}

func (r *CommandRepository) Delete(id string) error {
	return r.db.Delete(&db.CommandConfig{}, "id = ?", id).Error
}
