package command

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {
	route := r.Group("/command")

	route.POST("")
}
