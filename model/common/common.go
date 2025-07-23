package common

import "time"

/*
通用的结构体
*/

// Token token
type Token struct {
	Token    string    `json:"token,omitempty"` // token
	ExpireAt time.Time `json:"expire_at"`       // token 过期时间
}

// Page 分页
type Page struct {
	Page     int32 `json:"page" form:"page"`           // 第几页
	PageSize int32 `json:"page_size" form:"page_size"` // 每页的大小
}
