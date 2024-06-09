package types

import "time"

type Scheduler struct {
	Id       string    `json:"id"`
	Label    string    `json:"label"`
	LastPing time.Time `json:"lastPing" time_format:"unix"`
}
