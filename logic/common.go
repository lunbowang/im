package logic

import (
	"im/global"
	"im/model"
	"time"

	"github.com/XYYSWK/Lutils/pkg/token"
)

// newUserToken
// 成功：返回 token
func newUserToken(t model.TokenType, id int64, expireTime time.Duration) (string, *token.Payload, error) {
	if t == model.AccountToken {
		return "", nil, nil
	}
	// 设置令牌有效期
	duration := expireTime
	// 构建令牌内容并序列化为字节
	// NewTokenContent创建包含令牌类型和ID的内容对象
	// Marshal将内容序列化为JSON或其他格式
	data, err := model.NewTokenContent(t, id).Marshal()
	if err != nil {
		return "", nil, err
	}
	// 使用全局令牌生成器创建令牌
	result, payload, err := global.TokenMaker.CreateToken(data, duration)
	if err != nil {
		return "", nil, err
	}
	return result, payload, nil
}
