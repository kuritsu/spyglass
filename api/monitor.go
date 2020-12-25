package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/types"
)

// Monitors actions
type Monitors struct{}

// Get all monitors
func (m Monitors) Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"hello": "you",
	})
}

// Post a new monitor
func (m Monitors) Post(c *gin.Context) {
	var monitor types.Monitor
	if er := c.ShouldBind(&monitor); er == nil {
		log.Println(monitor.ID)
		log.Println(monitor.Schedule)
		log.Println(monitor.Type)
		log.Println(monitor.Definition)
		log.Println(monitor.Definition.DockerDefinition.DockerEnv)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": er.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id": monitor.ID,
	})
}
