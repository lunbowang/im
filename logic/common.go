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
	duration := expireTime
	data, err := model.NewTokenContent(t, id).Marshal()
	if err != nil {
		return "", nil, err
	}
	result, payload, err := global.TokenMaker.CreateToken(data, duration)
	if err != nil {
		return "", nil, err
	}
	return result, payload, nil
}
