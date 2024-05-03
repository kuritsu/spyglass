package types

type Role struct {
	Name        string `json:"name" yaml:"name" binding:"required"`
	Description string `json:"description" yaml:"description" binding:"required"`
	Permissions
}
