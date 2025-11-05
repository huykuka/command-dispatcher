package models

type CommandCreateDTO struct {
	Description string              `json:"description"`
	DeviceID    string              `json:"deviceId" validate:"required"`
	Type        string              `json:"type" validate:"required"`
	Parameters  []map[string]string `json:"parameters"`
}
