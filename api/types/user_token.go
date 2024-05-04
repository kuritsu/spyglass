package types

import "time"

type UserToken struct {
	Email      string    `json:"email" binding:"required"`
	Token      string    `json:"token" binding:"required" bson:"token"`
	Expiration time.Time `json:"expiration" binding:"required"`
}
