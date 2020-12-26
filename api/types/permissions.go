package types

import "time"

// Permissions for all objects
type Permissions struct {
	Owner     string    `json:"owner"`
	Readers   []string  `json:"readers"`
	Writers   []string  `json:"writers"`
	CreatedAt time.Time `json:"createdAt" time_format:"unix"`
	UpdatedAt time.Time `json:"updatedAt" time_format:"unix"`
}
