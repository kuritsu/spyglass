package types

// View for targets
type View struct {
	ImageBig   string `json:"imageBig"`
	ImageSmall string `json:"imageSmall"`
	ColorBig   string `json:"colorBig"`
	ColorSmall string `json:"colorSmall"`
}

// MonitorRef is a reference to an existing monitor
type MonitorRef struct {
	MonitorID string      `json:"monitorId"`
	Params    interface{} `json:"params"`
}

// Target full object
type Target struct {
	ID                string     `json:"id" binding:"required"`
	Description       string     `json:"description" binding:"required"`
	URL               string     `json:"url"`
	View              View       `json:"view" binding:"required"`
	Status            int        `json:"status"`
	StatusDescription string     `json:"statusDescription"`
	Critical          bool       `json:"critical"`
	Monitor           MonitorRef `json:"monitor"`
	Permissions
}
