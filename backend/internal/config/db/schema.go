package db

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        string     `json:"id" gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}

// CommandConfig defines the configuration for a specific command.
type CommandConfig struct {
	Base
	Name                  string `json:"name" gorm:"unique;not null"`
	Description           string `json:"description"`
	CommandType           string `json:"commandType" gorm:"not null"` // e.g., "rpc", "deviceData", "configuration"
	IsAcknowledgeRequired bool   `json:"isAcknowledgeRequired" gorm:"default:false"`
	PayloadSchema         string `json:"payloadSchema" gorm:"default:'{}'"` // JSON schema for validating command arguments/payload
	AcknowlegmentTimeout  int    `json:"acknowledgementTimeout" gorm:"default:60"`
	CompletionTimeout     int    `json:"completionTimeout" gorm:"default:60"`
}

// CommandExecution records the history and status of a command sent to a device.
type CommandExecution struct {
	Base
	DeviceID             string          `json:"deviceId" gorm:"not null;index"`
	CommandConfigID      string          `json:"commandConfigId" gorm:"type:uuid;not null"`
	CommandConfig        CommandConfig   `json:"-" gorm:"foreignKey:CommandConfigID"` // Belongs-to relationship
	Status               string          `json:"status" gorm:"index"`                 // e.g., "PENDING", "SENT", "ACKNOWLEDGED", "COMPLETED", "FAILED"
	IssuedAt             time.Time       `json:"issuedAt" gorm:"autoCreateTime"`
	CompletedAt          *time.Time      `json:"completedAt"`
	ExecutionHistory     json.RawMessage `json:"executionHistory" gorm:"type:jsonb"` // Store execution events as JSON
	CommandExecutionTime time.Time       `json:"commandExecutionTime"`
}
