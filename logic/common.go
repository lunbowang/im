package logic

import (
	"im/global"
	"im/model"
	"im/pkg/retry"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/XYYSWK/Lutils/pkg/token"
)

// 尝试重试
// 失败：打印日志
func reTry(name string, f func() error) {
	go func() {
		report := <-retry.NewTry(name, f, global.PublicSetting.Auto.Retry.Duration, global.PublicSetting.Auto.Retry.MaxTimes).Run()
		global.Logger.Error(report.Error())
	}()
}

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

// GetTokenAndPayload 获取 token 和 Payload
func GetTokenAndPayload(ctx *gin.Context) (string, *token.Payload, error) {
	// 从请求头中获取Authorization字段的值
	tokenString := ctx.GetHeader(global.PrivateSetting.Token.AuthorizationKey)
	result := strings.Fields(tokenString)
	tokenString = result[1]

	// 使用全局令牌生成器验证令牌有效性
	// 验证过程包括检查签名完整性、令牌是否过期等
	payload, err := global.TokenMaker.VerifyToken(tokenString)
	if err != nil {
		// 若验证失败，返回错误（如无效签名、令牌过期等）
		return "", nil, err
	}

	// 验证成功，返回令牌字符串和解析后的载荷
	return tokenString, payload, nil
}

// newAccountToken token
// 成功：返回 token，*token.Payload
// 失败：返回 nil, error
func newAccountToken(t model.TokenType, id int64) (string, *token.Payload, error) {
	if t != model.AccountToken {
		return "", nil, nil
	}
	duration := global.PrivateSetting.Token.AccountTokenDuration
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
