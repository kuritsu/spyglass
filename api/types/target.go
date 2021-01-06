package types

import "strings"

// View for targets
type View struct {
	ImageBig   string `json:"imageBig,omitempty" hcl:"image_big,optional"`
	ImageSmall string `json:"imageSmall,omitempty" hcl:"image_small,optional"`
	ColorBig   string `json:"colorBig,omitempty" hcl:"color_big,optional"`
	ColorSmall string `json:"colorSmall,omitempty" hcl:"color_small,optional"`
}

// MonitorRef is a reference to an existing monitor
type MonitorRef struct {
	MonitorID string            `json:"monitorId,omitempty" hcl:"monitor_id"`
	Params    map[string]string `json:"params,omitempty" bson:",omitempty" hcl:"params,optional"`
}

// Target full object
type Target struct {
	ID                string      `json:"id" binding:"required" hcl:"id,label"`
	Description       string      `json:"description" binding:"required" hcl:"description"`
	URL               string      `json:"url,omitempty" bson:",omitempty" hcl:"url,optional"`
	View              *View       `json:"view,omitempty" bson:",omitempty" hcl:"view,block"`
	Status            int         `json:"status" hcl:"status,optional"`
	StatusDescription string      `json:"statusDescription,omitempty" bson:",omitempty" hcl:"status_description,optional"`
	StatusTotal       int         `json:"-" bson:"statusTotal"`
	Critical          bool        `json:"critical" hcl:"critical,optional"`
	Monitor           *MonitorRef `json:"monitor,omitempty" bson:",omitempty" hcl:"monitor,block"`
	Children          []Target    `json:"children,omitempty" bson:",omitempty"`
	ChildrenCount     int         `json:"childrenCount" bson:"childrenCount"`
	Permissions       `hcl:",remain"`
}

// TargetPatch represents the fields that can be patched
type TargetPatch struct {
	Status            int    `json:"status" binding:"required"`
	StatusDescription string `json:"statusDescription"`
}

// GetTargetParentByID obtains the parent ID of a target given its ID
func GetTargetParentByID(id string) string {
	parts := strings.Split(id, ".")
	if len(parts) == 1 {
		return ""
	}
	return strings.Join(parts[0:len(parts)-1], ".")
}

// GetIDForRegex escapes special chars in the id for regex usage.
func GetIDForRegex(id string) string {
	id = strings.ToLower(id)
	return strings.ReplaceAll(strings.ReplaceAll(id, ".", `\.`), "-", `\-`)
}
