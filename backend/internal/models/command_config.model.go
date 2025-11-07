package models

import "command-dispatcher/internal/config/db"

type CommandConfigCreateDTO struct {
	Name                  string `json:"name,omitempty" validate:"required"`
	Description           string `json:"description,omitempty"`
	CommandType           string `json:"commandType" validate:"required"`
	IsAcknowledgeRequired bool   `json:"isAcknowledgeRequired,omitempty"`
	PayloadSchema         string `json:"payloadSchema,omitempty"`
	AcknowlegmentTimeout  int    `json:"acknowledgementTimeout,omitempty"`
	CompletionTimeout     int    `json:"completionTimeout,omitempty"`
}

// ToEntity converts DTO to database entity
func (dto *CommandConfigCreateDTO) ToEntity() *db.CommandConfig {
	return &db.CommandConfig{
		Name:                  dto.Name,
		Description:           dto.Description,
		CommandType:           dto.CommandType,
		IsAcknowledgeRequired: dto.IsAcknowledgeRequired,
		PayloadSchema:         dto.PayloadSchema,
		AcknowlegmentTimeout:  dto.AcknowlegmentTimeout,
		CompletionTimeout:     dto.CompletionTimeout,
	}
}

type CommandConfigUpdateDTO struct {
	Name                  *string `json:"name"`
	Description           *string `json:"description"`
	CommandType           *string `json:"commandType"`
	IsAcknowledgeRequired *bool   `json:"isAcknowledgeRequired"`
	PayloadSchema         *string `json:"payloadSchema"`
	AcknowlegmentTimeout  *int    `json:"acknowledgementTimeout"`
	CompletionTimeout     *int    `json:"completionTimeout"`
}

// ApplyTo safely updates entity with non-nil DTO fields
func (dto *CommandConfigUpdateDTO) ApplyTo(entity *db.CommandConfig) {
	if dto.Name != nil {
		entity.Name = *dto.Name
	}
	if dto.Description != nil {
		entity.Description = *dto.Description
	}
	if dto.CommandType != nil {
		entity.CommandType = *dto.CommandType
	}
	if dto.IsAcknowledgeRequired != nil {
		entity.IsAcknowledgeRequired = *dto.IsAcknowledgeRequired
	}
	if dto.PayloadSchema != nil {
		entity.PayloadSchema = *dto.PayloadSchema
	}
	if dto.AcknowlegmentTimeout != nil {
		entity.AcknowlegmentTimeout = *dto.AcknowlegmentTimeout
	}
	if dto.CompletionTimeout != nil {
		entity.CompletionTimeout = *dto.CompletionTimeout
	}
}
