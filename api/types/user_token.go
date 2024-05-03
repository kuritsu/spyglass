package types

import "time"

type UserToken struct {
	Email      string    `json:"email" binding:"required"`
	TokenHash  string    `json:"tokenHash" binding:"required" bson:"tokenHash"`
	Expiration time.Time `json:"expiration" binding:"required"`
}
