package middlewares

import (
	"github.com/gin-gonic/gin"
)

func NoJsonAPI() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("IsJsonAPI", false)
		c.Next()
	}
}
