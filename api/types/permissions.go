package types

// Permissions for all objects
type Permissions struct {
	Owner   string   `json:"owner"`
	Readers []string `json:"readers"`
	Writers []string `json:"writers"`
}
