package api

import "github.com/gin-gonic/gin"

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
}

// Monitor is a monitor definition for a target
type Monitor struct {
	Id         string            `json:"id"`
	Type       string            `json:"type"`
	Schedule   string            `json:"schedule"`
	Definition MonitorDefinition `json:"definition"`
}

// Monitors actions
type Monitors struct{}

// Get all monitors
func (m Monitors) Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"hello": "you",
	})
}
