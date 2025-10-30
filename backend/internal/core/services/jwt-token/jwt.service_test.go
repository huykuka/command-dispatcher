package jwttoken

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestNewJWTService(t *testing.T) {
	// Set up
	os.Setenv("HASH_JWT_KEY", "test_key")
	defer os.Unsetenv("HASH_JWT_KEY")

	// Test
	service1 := NewJWTService()
	service2 := NewJWTService()

	// Assert
	assert.NotNil(t, service1)
	assert.Equal(t, service1, service2, "NewJWTService should return the same instance")
	assert.Equal(t, []byte("test_key"), service1.JWTKey)
}

func TestGenerateAccessToken(t *testing.T) {
	// Set up
	os.Setenv("HASH_JWT_KEY", "test_key")
	defer os.Unsetenv("HASH_JWT_KEY")
	service := NewJWTService()

	// Test
	token, err := service.GenerateAccessToken("test@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test_key"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, "test@example.com", claims["username"])
	assert.NotNil(t, claims["exp"])
}

func TestValidateToken(t *testing.T) {
	// Set up
	os.Setenv("HASH_JWT_KEY", "test_key")
	defer os.Unsetenv("HASH_JWT_KEY")
	service := NewJWTService()

	// Generate a valid token
	validToken, _ := service.GenerateAccessToken("test@example.com")

	// Test cases
	testCases := []struct {
		name          string
		token         string
		expectedValid bool
	}{
		{"Valid token", validToken, true},
		{"Invalid token", "invalid.token.here", false},
		{"Empty token", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parsedToken, err := service.ValidateToken(tc.token)

			if tc.expectedValid {
				assert.NoError(t, err)
				assert.NotNil(t, parsedToken)
				assert.True(t, parsedToken.Valid)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	// Set up
	os.Setenv("HASH_JWT_KEY", "test_key")
	defer os.Unsetenv("HASH_JWT_KEY")
	service := NewJWTService()

	// Generate a token that expires in 1 second
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": "test@example.com",
			"exp":      time.Now().Add(time.Second).Unix(),
		})

	tokenString, _ := token.SignedString([]byte("test_key"))

	// Wait for the token to expire
	time.Sleep(2 * time.Second)

	// Validate the expired token
	_, err := service.ValidateToken(tokenString)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Token is expired")
}
