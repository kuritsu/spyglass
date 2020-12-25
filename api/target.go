package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

// View for targets
type View struct {
	ImageBig   string `json:"imageBig"`
	ImageSmall string `json:"imageSmall"`
	ColorBig   string `json:"colorBig"`
	ColorSmall string `json:"colorSmall"`
}

// MonitorRef is a reference to an existing monitor
type MonitorRef struct {
	MonitorID string      `json:"monitorId"`
	Params    interface{} `json:"params"`
}

// Target full object
type Target struct {
	ID                string     `json:"id" binding:"required"`
	Description       string     `json:"description" binding:"required"`
	URL               string     `json:"url"`
	View              View       `json:"view" binding:"required"`
	Status            int        `json:"status"`
	StatusDescription string     `json:"statusDescription"`
	Critical          bool       `json:"critical"`
	Monitor           MonitorRef `json:"monitor"`
	CreatedAt         time.Time  `json:"createdAt" time_format:"unix"`
	UpdatedAt         time.Time  `json:"updatedAt" time_format:"unix"`
}

// Targets actions
type Targets struct{}

// Get all targets
func (t Targets) Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
