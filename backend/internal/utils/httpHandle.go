package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Page struct {
	Size   int `json:"fromDate,omitempty" form:"page[size]" validate:"omitempty,gt=0"`
	Number int `json:"toDate,omitempty" form:"page[number]" validate:"omitempty,gt=0"`
}

type Pageable interface {
	GetPage() Page
}

// HandleHTTPError handles HTTP errors by logging the error message, setting error details in the context, and aborting the request with a status code.
//
// Parameters:
// - c: *gin.Context - The Gin context for the current request.
// - errMsg: string - The error message to be logged.
// - details: string - Additional details about the error to be set in the context.
// - statusCode: ...int - Optional variadic parameter for the HTTP status code. If not provided, defaults to 400 (Bad Request).
func HandleHTTPError(c *gin.Context, errMsg string, details string, statusCode ...int) {
	code := http.StatusBadRequest
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	log.WithFields(log.Fields{}).Error(errMsg)
	c.Set("error", details)
	c.Set("statusCode", strconv.Itoa(code))
	c.AbortWithStatus(code)
}

// SetResponse sets the response data in the context.
// Parameters:
// - c: *gin.Context - The Gin context for the current request.
// - data: map[string]any - The response data to be set in the context.

func SetResponse(c *gin.Context, data map[string]any) {
	c.Set("response", gin.H(data))
}

// CreatePaging creates a paging query for the provided GORM query and Pageable object.
// Parameters:
// - qr: *gorm.DB - The GORM query to be paginated.
// - query: T - The Pageable object containing the paging information.
// Example:
//
//	type GetSettingQueryDTO struct {
//		Filter SettingFilterDTO `json:"filter"`
//		Page   Page             `json:"page"`
func CreatePaging[T Pageable](qr *gorm.DB, query T) {
	page := query.GetPage()
	if page.Size > 0 {
		if page.Number > 0 {
			offset := (page.Number - 1) * page.Size
			qr = qr.Offset(offset).Limit(page.Size)
		} else {
			qr = qr.Limit(page.Size)
		}
	}
}
