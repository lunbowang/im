package server

/*
chat 中 server 端有关消息请求的结构
*/

type ReadMsg struct {
	EnToken  string  `json:"en_token,omitempty"` // 加密后的 Token
	MsgIDs   []int64 `json:"msg_ids"`            // 已读消息 IDs
	ReaderID int64   `json:"reader_id"`          // 读者账号 ID
}
