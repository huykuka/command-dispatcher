package models

type CommandConfigCreateDTO struct {
	Name                   string `json:"name,omitempty" validate:"required"`
	Description            string `json:"description,omitempty"`
	CommandType            string `json:"commandType" validate:"required"`
	IsAcknowledgeRequired  bool   `json:"isAcknowledgeRequired,omitempty"`
	PayloadSchema          string `json:"payloadSchema,omitempty"`
	AknowlegmentTimeout    int    `json:"acknowledgementTimeout,omitempty"`
	CompletionTimeout      int    `json:"completionTimeout,omitempty"`
}

type CommandConfigUpdateDTO struct {
	Name                   *string `json:"name"`
	Description            *string `json:"description"`
	CommandType            *string `json:"commandType"`
	IsAcknowledgeRequired  *bool   `json:"isAcknowledgeRequired"`
	PayloadSchema          *string `json:"payloadSchema"`
	AknowlegmentTimeout    *int    `json:"acknowledgementTimeout"`
	CompletionTimeout      *int    `json:"completionTimeout"`
}
