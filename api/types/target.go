package types

import "strings"

// View for targets
type View struct {
	ImageBig   string `json:"imageBig,omitempty"`
	ImageSmall string `json:"imageSmall,omitempty"`
	ColorBig   string `json:"colorBig,omitempty"`
	ColorSmall string `json:"colorSmall,omitempty"`
}

// MonitorRef is a reference to an existing monitor
type MonitorRef struct {
	MonitorID string            `json:"monitorId,omitempty"`
	Params    map[string]string `json:"params,omitempty" bson:",omitempty"`
}

// Target full object
type Target struct {
	ID                string      `json:"id" binding:"required"`
	Description       string      `json:"description" binding:"required"`
	URL               string      `json:"url,omitempty" bson:",omitempty"`
	View              *View       `json:"view,omitempty" bson:",omitempty"`
	Status            int         `json:"status"`
	StatusDescription string      `json:"statusDescription,omitempty" bson:",omitempty"`
	Critical          bool        `json:"critical"`
	Monitor           *MonitorRef `json:"monitor,omitempty" bson:",omitempty"`
	Permissions
}

// GetTargetParentByID obtains the parent ID of a target given its ID
func GetTargetParentByID(id string) string {
	parts := strings.Split(id, ".")
	if len(parts) == 1 {
		return ""
	}
	return strings.Join(parts[0:len(parts)-1], ".")
}
