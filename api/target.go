package api

import (
	"github.com/gin-gonic/gin"
)

// Targets actions
type Targets struct{}

// Get all targets
func (t Targets) Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
