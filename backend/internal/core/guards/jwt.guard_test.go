package guards_test

import (
	jwttoken "command-dispatcher/internal/core/services/jwt-token"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert" // Using a robust testing library

	"command-dispatcher/internal/core/guards"
)

func TestJWTAuthGuard_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", guards.JWTAuthGuard(), func(c *gin.Context) {
		c.String(http.StatusOK, "Success")
	})

	// Mock a valid token
	validToken, _ := jwttoken.NewJWTService().GenerateAccessToken("test@example.com")

	req, err := http.NewRequest(http.MethodGet, "/protected", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+validToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", w.Body.String())
}

func TestJWTAuthGuard_MissingAuthorizationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", guards.JWTAuthGuard(), func(c *gin.Context) {
		t.Errorf("Handler should not be called for missing authorization header")
	})
	req, err := http.NewRequest(http.MethodGet, "/protected", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuthGuard_PublicAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", guards.JWTAuthGuard(), func(c *gin.Context) {
		c.String(http.StatusOK, "Success")
	})

	req, err := http.NewRequest(http.MethodGet, "/protected", nil)
	assert.NoError(t, err)
	req.Header.Set("isPublic", "true")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", w.Body.String())
}

func TestJWTAuthGuard_InvalidAuthorizationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", guards.JWTAuthGuard(), func(c *gin.Context) {
		t.Errorf("Handler should not be called for invalid authorization header")
	})
	req, err := http.NewRequest(http.MethodGet, "/protected", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer invalid_token_string")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
