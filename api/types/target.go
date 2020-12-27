package types

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
	Params    map[string]string `json:"params,omitempty"`
}

// Target full object
type Target struct {
	ID                string     `json:"id" binding:"required"`
	Description       string     `json:"description" binding:"required"`
	URL               string     `json:"url,omitempty"`
	View              View       `json:"view" binding:"required"`
	Status            int        `json:"status"`
	StatusDescription string     `json:"statusDescription,omitempty"`
	Critical          bool       `json:"critical"`
	Monitor           MonitorRef `json:"monitor,omitempty"`
	Permissions
}
