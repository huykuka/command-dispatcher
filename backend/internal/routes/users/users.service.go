package users

import (
	"github.com/gin-gonic/gin"
)

type UserService struct{}

func (u *UserService) getAll(c *gin.Context) {
	c.Set("response", gin.H{
		"message": "pong",
	})
}

func (u *UserService) getByID(c *gin.Context) {
	id := c.Param("id")
	// Set the response data in the context
	c.Set("response", gin.H{
		"id":   id,
		"type": "user",
	})
}

func (u *UserService) create(c *gin.Context) {
	c.Set("response", gin.H{
		"message": "User created successfully",
	})
}

func (u *UserService) update(c *gin.Context) {
	id := c.Param("id")
	c.Set("response", gin.H{
		"message": "User updated successfully",
		"id":      id,
	})
}
