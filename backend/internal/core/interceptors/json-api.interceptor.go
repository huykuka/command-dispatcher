package interceptors

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type JSONAPIResponse struct {
	JSONAPI struct {
		Version string `json:"version"`
	} `json:"jsonapi"`
	Meta  map[string]interface{} `json:"meta,omitempty"`
	Links map[string]string      `json:"links,omitempty"`
	Data  interface{}            `json:"data,omitempty"`
	Error *JSONAPIError          `json:"error,omitempty"`
}

type JSONAPIError struct {
	Status string `json:"status,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func JsonApiInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		isJsonAPI, exist := c.Get("IsJsonAPI")
		if exist && !isJsonAPI.(bool) {
			// Bypass the middleware if IsJsonAPI is false
			c.Next()
			return
		}

		// Initialize the JSON API response
		jsonApiResponse := JSONAPIResponse{
			Meta:  make(map[string]interface{}),
			Links: make(map[string]string),
		}
		jsonApiResponse.JSONAPI.Version = "1.1"
		jsonApiResponse.Links["self"] = c.Request.URL.String()
		err, errorExist := c.Get("error")
		// Handle errors
		if errorExist {
			var statusCode int
			// Only append the first error from the Gin context to the JSON:API response
			jsonApiResponse.Error, statusCode = handleError(c, err.(string))
			// Set the status code and send the JSON:API formatted error response
			c.JSON(statusCode, jsonApiResponse)
			return
		}

		metadataResponse := handleGetMetadata(c)
		if metadataResponse != nil {
			jsonApiResponse.Meta = metadataResponse.Meta
		}

		// Get the actual response data
		responseData := handleResponseData(c)
		if responseData != nil {
			jsonApiResponse.Data = responseData.Data
		}
		// Replace the original response with the JSON-API formatted response
		c.JSON(c.Writer.Status(), jsonApiResponse)
	}
}

func handleResponseData(c *gin.Context) *JSONAPIResponse {
	jsonApiResponse := JSONAPIResponse{
		Meta: make(map[string]interface{}),
	}
	// Get the actual response data
	responseData, exists := c.Get("response")
	if exists {
		jsonApiResponse.Data = responseData
		if dataList, ok := responseData.([]interface{}); ok {
			jsonApiResponse.Meta["total"] = len(dataList)
		}
	}

	return &jsonApiResponse
}

func handleError(c *gin.Context, err string) (*JSONAPIError, int) {
	// Retrieve the status code from the context, defaulting to 400 if not set
	statusCode := http.StatusBadRequest

	if code, exists := c.Get("statusCode"); exists {
		statusCode, _ = strconv.Atoi(code.(string))
	}

	// Only append the first error from the Gin context to the JSON:API response
	return &JSONAPIError{
		Status: http.StatusText(statusCode),
		Title:  err, // Append only the first error
	}, statusCode
}

func handleGetMetadata(c *gin.Context) *JSONAPIResponse {

	jsonApiResponse := JSONAPIResponse{
		Meta: make(map[string]interface{}),
	}
	// Handle successful responses for GET requests
	if c.Request.Method == http.MethodGet {
		if page, exists := c.GetQuery("page[number]"); exists {
			jsonApiResponse.Meta["page"], _ = strconv.Atoi(page)
		}
		if size, exists := c.GetQuery("page[size]"); exists {
			jsonApiResponse.Meta["size"], _ = strconv.Atoi(size)
		}
		if sortField, exists := c.GetQuery("sort[field]"); exists {
			jsonApiResponse.Meta["sort"] = sortField
		}
		if sortOrder, exists := c.GetQuery("sort[order]"); exists {
			jsonApiResponse.Meta["order"] = sortOrder
		}
	}

	return &jsonApiResponse
}
