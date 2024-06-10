package types

// DockerDefinition is a Docker task definition
type DockerJobDefinition struct {
	Image   string            `json:"image,omitempty" bson:",omitempty" yaml:"image"`
	Command string            `json:"command,omitempty" bson:",omitempty" yaml:"command"`
	Env     map[string]string `json:"env,omitempty" bson:",omitempty" yaml:"env"`
	Volumes map[string]string `json:"volumes,omitempty" bson:",omitempty" yaml:"volumes"`
}

// K8SDefinition is a Kubernetes task definition
type K8SJobDefinition struct {
	Pod *DockerJobDefinition `json:"pod,omitempty" bson:",omitempty" yaml:"pod"`
}

// AWSServerlessDefinition is an AWS Lambda Serverless definition
type AWSJobDefinition struct {
	LambdaArn string `json:"lambdaArn,omitempty" bson:",omitempty" yaml:"lambda_arn"`
	Event     string `json:"event,omitempty" bson:",omitempty" yaml:"event"`
}

// ShellDefinition is a Shell command definition
type ShellJobDefinition struct {
	Command    string            `json:"command,omitempty" bson:",omitempty" yaml:"command"`
	Executable string            `json:"executable,omitempty" bson:",omitempty" yaml:"executable"`
	Env        map[string]string `json:"env,omitempty" bson:",omitempty" yaml:"env"`
}

// MonitorJobDefinition is a definition of a monitor
type MonitorJobDefinition struct {
	Docker *DockerJobDefinition `json:"docker,omitempty" bson:",omitempty" yaml:"docker"`
	K8S    *K8SJobDefinition    `json:"k8s,omitempty" bson:",omitempty" yaml:"k8s"`
	AWS    *AWSJobDefinition    `json:"aws,omitempty" bson:",omitempty" yaml:"aws"`
	Shell  *ShellJobDefinition  `json:"shell,omitempty" bson:",omitempty" yaml:"shell"`
}

// Monitor is a monitor definition to be assigned to targets
type Monitor struct {
	ID          string                `json:"id" binding:"required" yaml:"id"`
	Label       string                `json:"label" binding:"required" yaml:"label"`
	Schedule    string                `json:"schedule" binding:"required" yaml:"schedule"`
	Description string                `json:"description,omitempty" bson:",omitempty" yaml:"description"`
	Definition  *MonitorJobDefinition `json:"definition" binding:"required" yaml:"definition"`
	Permissions `yaml:",inline"`
}
