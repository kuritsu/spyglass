package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/api/types"
)

// MonitorController actions
type MonitorController struct {
	db storage.Provider
}

// Initialize the controller
func (m *MonitorController) Initialize(db storage.Provider) {
	m.db = db
}

// Get all monitors
func (m *MonitorController) Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"hello": "you",
	})
}

// Post a new monitor
func (m *MonitorController) Post(c *gin.Context) {
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
	m.db.Initialize()
	c.JSON(http.StatusCreated, gin.H{
		"id": monitor.ID,
	})
}
