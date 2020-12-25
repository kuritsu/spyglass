package types

// TargetDefinition is the definition of a target
type TargetDefinition struct {
	TargetID string `json:"targetId"`
}

// DockerDefinition is a Docker task definition
type DockerDefinition struct {
	Image      string            `json:"image"`
	Entrypoint string            `json:"entrypoint"`
	DockerEnv  map[string]string `json:"dockerEnv"`
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
	Command string            `json:"command"`
	Env     map[string]string `json:"env"`
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
	Permissions
}
