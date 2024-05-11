package types

type UserUpdateRequest struct {
	FullName    string `json:"fullName,omitempty" yaml:"fullName,omitempty"`
	OldPassword string `json:"oldPassword,omitempty" yaml:"fullName,omitempty"`
	NewPassword string `json:"newPassword,omitempty" yaml:"password,omitempty"`
	Permissions
}
