package types

import "time"

type UserTokenRequest struct {
	Expiration time.Time `json:"expiration" time_format:"unix"`
}
