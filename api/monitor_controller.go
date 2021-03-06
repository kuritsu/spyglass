package api

import (
	"net/http"
	"strconv"
	"strings"

	logr "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/api/types"
)

// MonitorController actions
type MonitorController struct {
	db  storage.Provider
	log *logr.Logger
}

// Init -ialize the controller
func (m *MonitorController) Init(db storage.Provider, log *logr.Logger) {
	m.db = db
	m.log = log
}

// Get a monitor by its Id
func (m *MonitorController) Get(c *gin.Context) {
	id := c.Param("id")
	m.db.Init()
	defer m.db.Free()
	monitor, err := m.db.GetMonitorByID(id)
	switch {
	case err != nil:
		m.log.Error(err)
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

// GetAll monitors, paginated
func (m *MonitorController) GetAll(c *gin.Context) {
	m.db.Init()
	defer m.db.Free()
	pageSizeString := c.DefaultQuery("pageSize", "10")
	pageIndexString := c.DefaultQuery("pageIndex", "0")
	pageSize, err := strconv.ParseInt(pageSizeString, 10, 64)
	if err != nil || pageSize > 100 || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid page size.",
		})
		return
	}
	pageIndex, err := strconv.ParseInt(pageIndexString, 10, 64)
	if err != nil || pageIndex < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid page index.",
		})
		return
	}
	contains := c.Query("contains")
	if contains != "" && !IsValidIDFragment(contains) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid contains expression.",
		})
		return
	}
	monitor, err := m.db.GetAllMonitors(pageSize, pageIndex, contains)
	if err != nil {
		m.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid operation. Try again.",
		})
		return
	}
	c.JSON(http.StatusOK, []types.Monitor(monitor))
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
	if !IsValidID(monitor.ID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid monitor ID.",
		})
		return
	}
	m.db.Init()
	defer m.db.Free()
	_, err := m.db.InsertMonitor(&monitor)
	if err != nil {
		m.log.Error(err)
		if strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "Duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Duplicate monitor ID.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid operation. Try again.",
		})
		return
	}
	c.JSON(http.StatusCreated, monitor)
}

// Put an existing monitor.
func (m *MonitorController) Put(c *gin.Context) {
	var newMonitor types.Monitor
	if er := c.ShouldBind(&newMonitor); er != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": er.Error(),
		})
		return
	}
	if !IsValidID(newMonitor.ID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid monitor ID.",
		})
		return
	}
	m.db.Init()
	defer m.db.Free()
	oldMonitor, err := m.db.GetMonitorByID(newMonitor.ID)
	if err != nil {
		m.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid operation. Try again.",
		})
		return
	}
	if oldMonitor == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Monitor not found.",
		})
		return
	}
	m.db.Init()
	defer m.db.Free()
	updatedMonitor, err := m.db.UpdateMonitor(oldMonitor, &newMonitor)
	if err != nil {
		m.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid operation. Try again.",
		})
		return
	}
	c.JSON(http.StatusOK, updatedMonitor)
}
