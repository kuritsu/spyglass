package types

// TargetDefinition is the definition of a target
type TargetDefinition struct {
	TargetID string `json:"targetId,omitempty" bson:",omitempty" hcl:"target_id,optional"`
}

// DockerDefinition is a Docker task definition
type DockerDefinition struct {
	Image      string            `json:"image,omitempty" bson:",omitempty" hcl:"image,optional"`
	Entrypoint string            `json:"entrypoint,omitempty" bson:",omitempty" hcl:"entrypoint,optional"`
	DockerEnv  map[string]string `json:"dockerEnv,omitempty" bson:",omitempty" hcl:"docker_env,optional"`
}

// K8SDefinition is a Kubernetes task definition
type K8SDefinition struct {
	Pod *DockerDefinition `bson:",omitempty" hcl:"pod,block"`
}

// AWSServerlessDefinition is an AWS Lambda Serverless definition
type AWSServerlessDefinition struct {
	LambdaArn string `json:"lambdaArn,omitempty" bson:",omitempty" hcl:"lambda_arn,optional"`
	Event     string `json:"event,omitempty" bson:",omitempty" hcl:"event,optional"`
}

// AzureServerlessDefinition is an Azure Function Serverless definition
type AzureServerlessDefinition struct {
	AzureFunc string `json:"azureFunc,omitempty" bson:",omitempty" hcl:"azure_func,optional"`
	Body      string `json:"body,omitempty" bson:",omitempty" hcl:"body,optional"`
}

// ServerlessDefinition is a Serverless definition
type ServerlessDefinition struct {
	AWS   *AWSServerlessDefinition   `bson:",omitempty" hcl:"aws,block"`
	Azure *AzureServerlessDefinition `bson:",omitempty" hcl:"azure,block"`
}

// ShellDefinition is a Shell command definition
type ShellDefinition struct {
	Command string            `json:"command,omitempty" bson:",omitempty" hcl:"command,optional"`
	Env     map[string]string `json:"env,omitempty" bson:",omitempty" hcl:"env,optional"`
}

// MonitorDefinition is a definition of a monitor
type MonitorDefinition struct {
	Docker     *DockerDefinition     `bson:",omitempty" hcl:"docker,block"`
	K8S        *K8SDefinition        `bson:",omitempty" hcl:"k8s,block"`
	Serverless *ServerlessDefinition `bson:",omitempty" hcl:"serverless,block"`
	Shell      *ShellDefinition      `bson:",omitempty" hcl:"shell,block"`
	Target     *TargetDefinition     `bson:",omitempty" hcl:"target,block"`
}

// Monitor is a monitor definition for a target
type Monitor struct {
	ID          string `json:"id" binding:"required" hcl:"id,label"`
	Type        string `json:"type" binding:"required" hcl:"type"`
	Schedule    string `json:"schedule" binding:"required" hcl:"schedule"`
	Permissions `hcl:",remain"`
	Definition  *MonitorDefinition `json:"definition" binding:"required" hcl:"definition,block"`
}
