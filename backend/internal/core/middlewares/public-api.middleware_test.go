package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPublicApiMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Variable to capture the value set by the middleware
	var isPublicValue bool

	// Create a new Gin engine
	r := gin.New()

	// Use the middleware you want to test
	r.Use(PublicApiMiddleware())

	// Define a simple handler to capture the value of "isPublic" from the context
	r.GET("/", func(c *gin.Context) {
		isPublic, exists := c.Get("isPublic")
		if exists {
			isPublicValue = isPublic.(bool) // Capture the value from the context
		} else {
			isPublicValue = false // Default to false if not set
		}
		c.Status(http.StatusOK)
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	// Serve the request
	r.ServeHTTP(w, req)

	// Assert that the middleware set the "isPublic" value to false
	assert.Equal(t, true, isPublicValue)
	assert.Equal(t, http.StatusOK, w.Code)
}
