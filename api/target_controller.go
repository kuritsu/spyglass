package api

import (
	"github.com/gin-gonic/gin"
)

// TargetController actions
type TargetController struct{}

// Get all targets
func (t TargetController) Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
