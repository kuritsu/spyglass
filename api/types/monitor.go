package types

// TargetDefinition is the definition of a target
type TargetDefinition struct {
	TargetID string `json:"targetId,omitempty" bson:",omitempty" hcl:"target_id"`
}

// DockerDefinition is a Docker task definition
type DockerDefinition struct {
	Image      string            `json:"image,omitempty" bson:",omitempty" hcl:"image"`
	Entrypoint string            `json:"entrypoint,omitempty" bson:",omitempty" hcl:"entrypoint"`
	DockerEnv  map[string]string `json:"dockerEnv,omitempty" bson:",omitempty" hcl:"docker_env"`
}

// K8SDefinition is a Kubernetes task definition
type K8SDefinition struct {
	DockerDefinition `bson:",omitempty"`
}

// AWSServerlessDefinition is an AWS Lambda Serverless definition
type AWSServerlessDefinition struct {
	LambdaArn string `json:"lambdaArn,omitempty" bson:",omitempty" hcl:"lambda_arn"`
	Event     string `json:"event,omitempty" bson:",omitempty" hcl:"event"`
}

// AzureServerlessDefinition is an Azure Function Serverless definition
type AzureServerlessDefinition struct {
	AzureFunc string `json:"azureFunc,omitempty" bson:",omitempty" hcl:"azure_func"`
	Body      string `json:"body,omitempty" bson:",omitempty" hcl:"body"`
}

// ServerlessDefinition is a Serverless definition
type ServerlessDefinition struct {
	AWSServerlessDefinition   `bson:",omitempty"`
	AzureServerlessDefinition `bson:",omitempty"`
}

// ShellDefinition is a Shell command definition
type ShellDefinition struct {
	Command string            `json:"command,omitempty" bson:",omitempty" hcl:"command"`
	Env     map[string]string `json:"env,omitempty" bson:",omitempty" hcl:"env"`
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
	ID         string            `json:"id" binding:"required" hcl:"id,label"`
	Type       string            `json:"type" binding:"required" hcl:"type,label"`
	Schedule   string            `json:"schedule" binding:"required" hcl:"schedule"`
	Definition MonitorDefinition `json:"definition" binding:"required" hcl:"definition,block"`
	Permissions
}
