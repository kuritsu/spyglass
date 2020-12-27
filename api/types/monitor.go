package types

// TargetDefinition is the definition of a target
type TargetDefinition struct {
	TargetID string `json:"targetId,omitempty" bson:",omitempty"`
}

// DockerDefinition is a Docker task definition
type DockerDefinition struct {
	Image      string            `json:"image,omitempty" bson:",omitempty"`
	Entrypoint string            `json:"entrypoint,omitempty" bson:",omitempty"`
	DockerEnv  map[string]string `json:"dockerEnv,omitempty" bson:",omitempty"`
}

// K8SDefinition is a Kubernetes task definition
type K8SDefinition struct {
	DockerDefinition `bson:",omitempty"`
}

// AWSServerlessDefinition is an AWS Lambda Serverless definition
type AWSServerlessDefinition struct {
	LambdaArn string `json:"lambdaArn,omitempty" bson:",omitempty"`
	Event     string `json:"event,omitempty" bson:",omitempty"`
}

// AzureServerlessDefinition is an Azure Function Serverless definition
type AzureServerlessDefinition struct {
	AzureFunc string `json:"azureFunc,omitempty" bson:",omitempty"`
	Body      string `json:"body,omitempty" bson:",omitempty"`
}

// ServerlessDefinition is a Serverless definition
type ServerlessDefinition struct {
	AWSServerlessDefinition   `bson:",omitempty"`
	AzureServerlessDefinition `bson:",omitempty"`
}

// ShellDefinition is a Shell command definition
type ShellDefinition struct {
	Command string            `json:"command,omitempty" bson:",omitempty"`
	Env     map[string]string `json:"env,omitempty" bson:",omitempty"`
}

// MonitorDefinition is a definition of a monitor
type MonitorDefinition struct {
	DockerDefinition     `bson:",omitempty"`
	K8SDefinition        `bson:",omitempty"`
	ServerlessDefinition `bson:",omitempty"`
	ShellDefinition      `bson:",omitempty"`
	TargetDefinition     `bson:",omitempty"`
}

// Monitor is a monitor definition for a target
type Monitor struct {
	ID         string            `json:"id" binding:"required"`
	Type       string            `json:"type" binding:"required"`
	Schedule   string            `json:"schedule" binding:"required"`
	Definition MonitorDefinition `json:"definition" binding:"required"`
	Permissions
}
