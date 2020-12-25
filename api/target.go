package api

import "github.com/gin-gonic/gin"

// View for targets
type View struct {
	ImageBig   string `json:"imageBig"`
	ImageSmall string `json:"imageSmall"`
	ColorBig   string `json:"colorBig"`
	ColorSmall string `json:"colorSmall"`
}

// MonitorRef is a reference to an existing monitor
type MonitorRef struct {
	MonitorId string      `json:"monitorId"`
	Params    interface{} `json:"params"`
}

// Target full object
type Target struct {
	Id                string     `json:"id"`
	Description       string     `json:"description"`
	Url               string     `json:"url"`
	View              View       `json:"view"`
	Status            int        `json:"status"`
	StatusDescription string     `json:"statusDescription"`
	Critical          bool       `json:"critical"`
	Monitor           MonitorRef `json:"monitor"`
}

// Targets actions
type Targets struct{}

// Get all targets
func (t Targets) Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
