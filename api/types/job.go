package types

import "time"

type Job struct {
	ID          string            `json:"id"`
	TargetId    string            `json:"targetId"`
	Label       string            `json:"label"`
	Props       map[string]string `json:"props"`
	Status      int               `json:"status"`
	SchedulerId string            `json:"schedulerId"`
	UpdatedAt   time.Time         `json:"updatedAt" time_format:"unix"`
}
