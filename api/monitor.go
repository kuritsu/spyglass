package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TargetDefinition is the definition of a target
type TargetDefinition struct {
	TargetID string `json:"targetId"`
}

// DockerDefinition is a Docker task definition
type DockerDefinition struct {
	Image      string      `json:"image"`
	Entrypoint string      `json:"entrypoint"`
	Env        interface{} `json:"env"`
}

// K8SDefinition is a Kubernetes task definition
type K8SDefinition struct {
	DockerDefinition
}

// AWSServerlessDefinition is an AWS Lambda Serverless definition
type AWSServerlessDefinition struct {
	LambdaArn string `json:"lambdaArn"`
	Event     string `json:"event"`
}

// AzureServerlessDefinition is an Azure Function Serverless definition
type AzureServerlessDefinition struct {
	AzureFunc string `json:"azureFunc"`
	Body      string `json:"body"`
}

// ServerlessDefinition is a Serverless definition
type ServerlessDefinition struct {
	AWSServerlessDefinition
	AzureServerlessDefinition
}

// ShellDefinition is a Shell command definition
type ShellDefinition struct {
	Command string      `json:"command"`
	Env     interface{} `json:"env"`
}

// MonitorDefinition is a definition of a monitor
type MonitorDefinition struct {
	DockerDefinition
	K8SDefinition
	ServerlessDefinition
	ShellDefinition
	TargetDefinition
}

// Monitor is a monitor definition for a target
type Monitor struct {
	ID         string            `json:"id" binding:"required"`
	Type       string            `json:"type" binding:"required"`
	Schedule   string            `json:"schedule" binding:"required"`
	Definition MonitorDefinition `json:"definition" binding:"required"`
}

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
	var monitor Monitor
	if er := c.ShouldBind(&monitor); er == nil {
		log.Println(monitor.ID)
		log.Println(monitor.Schedule)
		log.Println(monitor.Type)
		log.Println(monitor.Definition)
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
