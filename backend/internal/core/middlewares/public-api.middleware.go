package middlewares

import "github.com/gin-gonic/gin"

func PublicApiMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("isPublic", true)
		c.Next()
	}
}
