package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a new rate limiter allowing 1 request per second
	limiter := rate.NewLimiter(1, 1)

	// Create a new Gin engine
	r := gin.New()
	r.Use(RateLimiter(limiter))

	// Define a simple handler
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Define test cases
	tests := []struct {
		name           string
		sleepDuration  time.Duration
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Should first request pass",
			sleepDuration:  0,
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			sleepDuration:  0,
			expectedStatus: http.StatusTooManyRequests,
			expectedBody:   "Rate limit exceed",
		},
		{
			name:           "should third request pass after waiting",
			sleepDuration:  time.Second,
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)

			// Wait if needed
			if tt.sleepDuration > 0 {
				time.Sleep(tt.sleepDuration)
			}

			// Serve the request
			r.ServeHTTP(w, req)

			// Check the response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
