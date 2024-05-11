package types

import "time"

// Permissions for all objects
type Permissions struct {
	Owners    []string  `json:"owners" yaml:"owners"`
	Readers   []string  `json:"readers" yaml:"readers"`
	Writers   []string  `json:"writers" yaml:"writers"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt" time_format:"unix"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt" time_format:"unix"`
}
