package model

import "time"

type SettingInfo struct {
	RelationID   int64     `json:"relation_id,omitempty"`    // 关系 ID
	RelationType string    `json:"relation_type,omitempty"`  // 关系类型['group','friend']
	NickName     string    `json:"nick_name,omitempty"`      // 昵称（群组时：账号在群中的昵称，好友时：给好友备注的昵称，空表示未设置）
	IsNotDisturb bool      `json:"is_not_disturb,omitempty"` // 是否免打扰
	IsPin        bool      `json:"is_pin,omitempty"`         // 是否 pin
	IsShow       bool      `json:"is_show,omitempty"`        // 是否显示
	PinTime      time.Time `json:"pin_time"`                 // pin 时间
	LastShow     time.Time `json:"last_show"`                // 最后显示时间
}

type SettingGroupInfo struct {
	RelationID  int64  `json:"relation_id,omitempty"` // 群组 ID
	Name        string `json:"name,omitempty"`        // 群组名称
	Description string `json:"description,omitempty"` // 群组的描述
	Avatar      string `json:"avatar,omitempty"`      // 群组头像
}

type SettingGroup struct {
	SettingInfo
	GroupInfo *SettingGroupInfo `json:"group_info"` // 群组信息
}
