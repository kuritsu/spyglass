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
	id := c.Param("id")
	m.db.Init()
	defer m.db.Free()
	monitor, err := m.db.GetMonitorByID(id)
	switch {
	case err != nil:
		log.Println("ERROR: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal error. Try again.",
		})
		return
	case monitor == nil:
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Monitor not found.",
		})
		return
	}
	c.JSON(http.StatusOK, monitor)
}

// Post a new monitor
func (m *MonitorController) Post(c *gin.Context) {
	var monitor types.Monitor
	if er := c.ShouldBind(&monitor); er != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": er.Error(),
		})
		return
	}
	m.db.Init()
	defer m.db.Free()
	_, err := m.db.InsertMonitor(&monitor)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid operation. Try again.",
		})
		return
	}
	c.JSON(http.StatusCreated, monitor)
}
