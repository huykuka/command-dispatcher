package pipes

import (
	"command-dispatcher/internal/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type SampleStruct struct {
	Name  string `json:"name" validate:"required"`
	Age   int    `json:"age" validate:"required,min=0,max=150"`
	Email string `json:"email" validate:"required,email"`
}

func TestStructValidator(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name      string
		input     SampleStruct
		expectErr bool
	}{
		{
			name: "Valid input",
			input: SampleStruct{
				Name:  "John Doe",
				Age:   140,
				Email: "john.doe@example.com",
			},
			expectErr: false,
		},
		{
			name: "Missing name",
			input: SampleStruct{
				Age:   30,
				Email: "john.doe@example.com",
			},
			expectErr: true,
		},
		{
			name: "Invalid email",
			input: SampleStruct{
				Name:  "John Doe",
				Age:   30,
				Email: "invalid-email",
			},
			expectErr: true,
		},
		{
			name: "Missing age",
			input: SampleStruct{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Call the structValidate function
			structValidate(c, &tt.input)

			// Check if an error was expected
			if tt.expectErr {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.NotEqual(t, http.StatusBadRequest, w.Code)
			}
		})
	}
}

func TestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validate = validator.New() // Ensure the validator is initialized

	tests := []struct {
		name      string
		body      string
		expectErr bool
	}{
		{
			name:      "Valid body",
			body:      `{"name":"John Doe","age":30,"email":"john.doe@example.com"}`,
			expectErr: false,
		},
		{
			name:      "Missing name",
			body:      `{"age":30,"email":"john.doe@example.com"}`,
			expectErr: true,
		},
		{
			name:      "Invalid email",
			body:      `{"name":"John Doe","age":30,"email":"invalid-email"}`,
			expectErr: true,
		},
		{
			name:      "Missing age",
			body:      `{"name":"John Doe","email":"john.doe@example.com"}`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the Body function
			Body[SampleStruct](c)

			// Check if an error was expected
			if tt.expectErr {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.NotEqual(t, http.StatusBadRequest, w.Code)
			}
		})
	}
}

type QueryBackupDTO struct {
	Filter struct {
		FromDate string `json:"fromDate,omitempty" form:"filter[fromDate]" validate:"omitempty"`
		ToDate   string `json:"toDate,omitempty" form:"filter[toDate]" validate:"omitempty"`
	} `json:"filter,omitempty" form:"filter"`
	utils.Page `json:"page,omitempty" form:"page"`
}

func TestQueryBackupDTO(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validate = validator.New() // Ensure the validator is initialized

	tests := []struct {
		name      string
		query     string
		expectErr bool
	}{
		{
			name:      "Valid query",
			query:     "filter[fromDate]=2023-09-03T14:00:00Z&filter[toDate]=2023-09-03T14:00:00Z&page[number]=1&page[size]=10",
			expectErr: false,
		},
		{
			name:      "Invalid date format",
			query:     "filter[fromDate]=invalid-date&filter[toDate]=2023-12-31&page[number]=1&page[size]=10",
			expectErr: false,
		},
		{
			name:      "Missing page number",
			query:     "filter[fromDate]=2023-01-01&filter[toDate]=2023-12-31&page[size]=10",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/?"+tt.query, nil)

			// Call the Query function
			Query[QueryBackupDTO](c)

			// Check if an error was expected
			if tt.expectErr {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.NotEqual(t, http.StatusBadRequest, w.Code)
			}
		})
	}
}
