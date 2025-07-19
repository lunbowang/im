package model

import "encoding/json"

type TokenType string

const (
	UserToken    TokenType = "user"
	AccountToken TokenType = "account"
)

type Content struct {
	TokenType TokenType `json:"token_type,omitempty"` // token 类型，用户token\账户token
	ID        int64     `json:"id,omitempty"`
}

// Unmarshal 将json序列解析为Content 结构体
func (c *Content) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	return nil
}

// NewTokenContent 新建一个类型的 token
func NewTokenContent(t TokenType, ID int64) *Content {
	return &Content{
		TokenType: t,
		ID:        ID,
	}
}

func (c *Content) Marshal() ([]byte, error) {
	return json.Marshal(c)
}
