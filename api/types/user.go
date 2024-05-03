package types

type User struct {
	Email     string `json:"email" yaml:"email" binding:"required"`
	FullName  string `json:"fullName" yaml:"fullName"`
	PassHash  string
	FirstHash string   `json:"firstHash"`
	Roles     []string `json:"roles,omitempty" yaml:"roles,omitempty"`
	Permissions
}
