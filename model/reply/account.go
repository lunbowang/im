package reply

import "im/model/common"

/*
定义账户相关的响应参数结构体
*/

type ParamCreateAccount struct {
	ParamAccountInfo     ParamAccountInfo     `json:"param_account_info"`      // 账号信息
	ParamGetAccountToken ParamGetAccountToken `json:"param_get_account_token"` // 账号 token
}

type ParamAccountInfo struct {
	ID     int64  `json:"id,omitempty"`     // 账号 ID
	Name   string `json:"name,omitempty"`   // 昵称
	Avatar string `json:"avatar,omitempty"` // 头像
	Gender string `json:"gender,omitempty"` // 性别 [男，女，未知]
}

type ParamGetAccountToken struct {
	AccountToken common.Token `json:"account_token"` // 账号 token
}
