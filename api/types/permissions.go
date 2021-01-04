package types

import "time"

// Permissions for all objects
type Permissions struct {
	Owner     string    `json:"owner"`
	Readers   []string  `json:"readers" hcl:"readers,optional"`
	Writers   []string  `json:"writers" hcl:"writers,optional"`
	CreatedAt time.Time `json:"createdAt" time_format:"unix"`
	UpdatedAt time.Time `json:"updatedAt" time_format:"unix"`
}
