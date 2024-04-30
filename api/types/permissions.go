package types

import "time"

// Permissions for all objects
type Permissions struct {
	Owners    []string  `json:"owners"`
	Readers   []string  `json:"readers"`
	Writers   []string  `json:"writers"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt" time_format:"unix"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt" time_format:"unix"`
}
