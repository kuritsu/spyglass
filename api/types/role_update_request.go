package types

type RoleUpdateRequest struct {
	UsersAdd    []string `json:"usersAdd,omitempty" yaml:"usersAdd,omitempty"`
	UsersRemove []string `json:"usersRemove,omitempty" yaml:"usersRemove,omitempty"`
	Permissions
}
