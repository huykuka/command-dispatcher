package interceptors

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name       string
		statusCode string
		error      string
	}{
		{
			name:       "Should return 401 error",
			statusCode: "401",
			error:      "Throttle limit exceed",
		},
		{
			name:       "Should return 404 error",
			statusCode: "404",
			error:      "Not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("error", tt.error)
			c.Set("statusCode", tt.statusCode)
			errorResponse, _ := handleError(c, tt.error)

			// Convert statusCode to int and then to http.StatusText
			statusCodeInt, err := strconv.Atoi(tt.statusCode)
			if err != nil {
				t.Fatalf("invalid status code: %v", err)
			}
			expectedStatusText := http.StatusText(statusCodeInt)

			assert.Equal(t, tt.error, errorResponse.Title)
			assert.Equal(t, expectedStatusText, errorResponse.Status)
		})
	}
}

func TestHandleResponseData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		responseData interface{}
		expectedData interface{}
		expectedMeta map[string]interface{}
	}{
		{
			name:         "Response data is a list",
			responseData: []interface{}{"item1", "item2", "item3"},
			expectedData: []interface{}{"item1", "item2", "item3"},
			expectedMeta: map[string]interface{}{"total": 3},
		},
		{
			name:         "Response data is a single item",
			responseData: "singleItem",
			expectedData: "singleItem",
			expectedMeta: map[string]interface{}{},
		},
		{
			name:         "No response data",
			responseData: nil,
			expectedData: nil,
			expectedMeta: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("response", tt.responseData)

			jsonApiResponse := handleResponseData(c)

			assert.Equal(t, tt.expectedData, jsonApiResponse.Data)
			assert.Equal(t, tt.expectedMeta, jsonApiResponse.Meta)
		})
	}
}

func TestHandleGetMetadata(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		query        string
		expectedMeta map[string]interface{}
	}{
		{
			name:         "All parameters present",
			query:        "?page[number]=1&page[size]=10&sort[field]=name&sort[order]=asc",
			expectedMeta: map[string]interface{}{"page": 1, "size": 10, "sort": "name", "order": "asc"},
		},
		{
			name:         "Missing sort parameters",
			query:        "?page[number]=2&page[size]=20",
			expectedMeta: map[string]interface{}{"page": 2, "size": 20},
		},
		{
			name:         "No parameters",
			query:        "",
			expectedMeta: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/"+tt.query, nil)

			jsonApiResponse := handleGetMetadata(c)

			assert.Equal(t, tt.expectedMeta, jsonApiResponse.Meta)
		})
	}
}

func TestJsonApiInterceptor(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		query          string
		setError       bool
		responseData   interface{}
		expectedStatus int
		expectedMeta   map[string]interface{}
		expectedData   interface{}
		expectedError  *JSONAPIError
	}{
		{
			name:           "Error handling",
			method:         http.MethodGet,
			query:          "",
			setError:       true,
			responseData:   nil,
			expectedStatus: http.StatusUnauthorized,
			expectedMeta:   map[string]interface{}(nil),
			expectedData:   nil,
			expectedError:  &JSONAPIError{Status: "Unauthorized", Title: "Throttle limit exceed"},
		},
		{
			name:           "POST request without metadata",
			method:         http.MethodPost,
			query:          "",
			setError:       false,
			responseData:   "singleItem",
			expectedStatus: http.StatusOK,
			expectedMeta:   map[string]interface{}(nil),
			expectedData:   "singleItem",
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(tt.method, "/"+tt.query, nil)

			if tt.setError {
				c.Set("error", "Throttle limit exceed")
				c.Set("statusCode", "401")
			}

			if tt.responseData != nil {
				c.Set("response", tt.responseData)
			}

			// Call the interceptor
			JsonApiInterceptor()(c)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse the JSON response
			var jsonApiResponse JSONAPIResponse
			err := json.Unmarshal(w.Body.Bytes(), &jsonApiResponse)
			assert.NoError(t, err)

			// Check the meta data
			assert.EqualValues(t, tt.expectedMeta, jsonApiResponse.Meta)

			// Check the response data
			assert.Equal(t, tt.expectedData, jsonApiResponse.Data)

			// Check the error
			assert.Equal(t, tt.expectedError, jsonApiResponse.Error)
		})
	}
}
