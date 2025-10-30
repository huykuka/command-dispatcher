package guards

import (
	jwttoken "command-dispatcher/internal/core/services/jwt-token"
	"command-dispatcher/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type JWTService interface {
	ValidateToken(tokenString string) (*jwt.Token, error)
}

func JWTAuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		isPublic := c.GetHeader("isPublic")
		if isPublic == "true" {
			c.Next()
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.HandleHTTPError(c, "Authorization header is required", "Permission denied", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			utils.HandleHTTPError(c, "Invalid authorization header format", "Permission denied", http.StatusUnauthorized)
			return
		}

		token, err := jwttoken.NewJWTService().ValidateToken(bearerToken[1])

		if err != nil || !token.Valid {
			utils.HandleHTTPError(c, "Invalid or expired token", "Permission denied", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.HandleHTTPError(c, "Failed to parse token claims", "Permission denied", http.StatusUnauthorized)
			return
		}

		// You can add the claims to the context if needed
		c.Set("username", claims["email"])

		c.Next()
	}
}
