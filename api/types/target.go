package types

import (
	"strings"
)

// MonitorRef is a reference to an existing monitor
type MonitorRef struct {
	MonitorID string            `json:"monitorId,omitempty" bson:"monitorId" hcl:"monitor_id"`
	Params    map[string]string `json:"params,omitempty" bson:",omitempty" hcl:"params,optional"`
}

// Target full object
type Target struct {
	ID                string      `json:"id" yaml:"id" binding:"required"`
	Description       string      `json:"description" yaml:"description" binding:"required"`
	URL               string      `json:"url,omitempty" yaml:"url,omitempty" bson:",omitempty"`
	Status            int         `json:"status" yaml:"status"`
	StatusDescription string      `json:"statusDescription,omitempty" yaml:"statusDescription,omitempty" bson:"statusDescription,omitempty"`
	Critical          bool        `json:"critical" yaml:"critical"`
	Monitor           *MonitorRef `json:"monitor,omitempty" bson:",omitempty"`
	Children          []TargetRef `json:"children,omitempty" yaml:"children,omitempty" bson:",omitempty"`
	ChildrenRef       []string    `json:"childrenRef,omitempty" yaml:"childrenRef,omitempty" bson:",omitempty"`
	Permissions
}

type TargetRef *Target

// TargetPatch represents the fields that can be patched
type TargetPatch struct {
	Status            int    `json:"status" binding:"required"`
	StatusDescription string `json:"statusDescription"`
}

// GetTargetParentByID obtains the parent ID of a target given its ID
func GetTargetParentByID(id string) string {
	id = strings.ToLower(id)
	parts := strings.Split(id, "/")
	if len(parts) == 1 {
		return ""
	}
	return strings.Join(parts[0:len(parts)-1], "/")
}

// GetShortID for children ref
func GetShortID(id string) string {
	parts := strings.Split(id, "/")
	if len(parts) == 1 {
		return id
	}
	return parts[len(parts)-1]
}

// GetIDForRegex escapes special chars in the id for regex usage.
func GetIDForRegex(id string) string {
	id = strings.ToLower(id)
	return strings.ReplaceAll(strings.ReplaceAll(id, ".", `\.`), "-", `\-`)
}

// TargetList for sorting targets by ID
type TargetList []*Target

func (s TargetList) Len() int {
	return len(s)
}

func (s TargetList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s TargetList) Less(i, j int) bool {
	return strings.ToLower(s[i].ID) < strings.ToLower(s[j].ID)
}
