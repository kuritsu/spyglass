package types

import "time"

type UserToken struct {
	TokenHash  string    `json:"tokenHash" binding:"required" bson:"tokenHash"`
	Expiration time.Time `json:"expiration" binding:"required"`
}
