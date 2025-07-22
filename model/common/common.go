package common

import "time"

// Token token
type Token struct {
	Token    string    `json:"token,omitempty"` // token
	ExpireAt time.Time `json:"expire_at"`       // token 过期时间
}
