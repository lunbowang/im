package middlewares

import (
	"im/dao"
	"im/errcodes"
	"im/global"
	"im/model"
	"net/http"
	"strings"

	"github.com/XYYSWK/Lutils/pkg/app"

	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/XYYSWK/Lutils/pkg/token"
	"github.com/gin-gonic/gin"
)

/*
用户验证（parseTo 生成 Token）
*/

// GetToken 从当前请求头获取 token
func GetToken(header http.Header) (string, errcode.Err) {
	// 本项目 Token 放在Header的Authorization 中，并使用 Bearer 开头
	authorizationHeader := header.Get(global.PrivateSetting.Token.AuthorizationKey)
	if len(authorizationHeader) == 0 {
		return "", errcodes.AuthNotExist
	}
	// 按空格切割（切割为：Bearer 和 token)
	parts := strings.SplitN(authorizationHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == global.PrivateSetting.Token.AuthorizationType) {
		return "", errcodes.AuthenticationFailed
	}
	return parts[1], nil
}

// ParseToken 解析header中的token。返回 payload,token,err
func ParseToken(accessToken string) (*token.Payload, string, errcode.Err) {
	// 解析 token
	payload, err := global.TokenMaker.VerifyToken(accessToken)
	if err != nil {
		if err.Error() == "超时错误" {
			return nil, "", errcodes.AuthOverTime
		}
		return nil, "", errcodes.AuthenticationFailed
	}
	return payload, accessToken, nil
}

// ParseToAuth 鉴权中间件，用于解析并写入 Token 信息入上下文中
// 该中间件会尝试从请求头中获取token，解析并验证有效性，然后将用户信息存入请求上下文中
// 后续的处理函数可以通过 ctx.Get(global.PrivateSetting.Token.AuthorizationKey) 获取当前用户信息
func ParseToAuth() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// 从请求头中提取token
		accessToken, err := GetToken(ctx.Request.Header)
		// 如果提取Token失败（Token 不存在），不做拦截直接继续处理请求
		if err != nil {
			ctx.Next()
			return
		}

		// 解析Token获取载荷信息
		// payload包含了Token中携带的用户身份信息
		payload, _, err := ParseToken(accessToken)
		// 如果Token解析失败（如Token过期、签名错误），继续处理请求但不注入用户信息
		if err != nil {
			ctx.Next()
			return
		}
		// 初始化Content结构体用于存储解析后的用户信息
		content := &model.Content{}
		// 将Token载荷中的Content字段反序列化为具体的Content对象
		if err := content.Unmarshal(payload.Content); err != nil {
			ctx.Next()
			return
		}
		// 将当前请求头中的 Content (token 类型和 id) 信息保存到请求的上下文 ctx 上
		ctx.Set(global.PrivateSetting.Token.AuthorizationKey, content)
		// 后续的处理请求的函数可以通过  ctx.Get(global.PrivateSetting.Token.AuthorizationKey)来获取当前请求的用户信息
		ctx.Next()
	}
}

// MustUser 必须是用户
func MustUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 创建响应处理器，统一返回错误信息
		reply := app.NewResponse(ctx)

		// 从上下文获取令牌数据
		val, ok := ctx.Get(global.PrivateSetting.Token.AuthorizationKey)
		if !ok {
			// 未找到令牌
			reply.Reply(errcodes.AuthNotExist)
			ctx.Abort()
			return
		}
		// 类型断言，确保令牌数据类型正确
		data := val.(*model.Content)

		if data.TokenType != model.UserToken {
			reply.Reply(errcodes.AuthenticationFailed)
			ctx.Abort()
			return
		}

		// 检查数据库中是否存在该用户（通过令牌中的用户ID查询）
		ok, err := dao.Database.DB.ExistsUserByID(ctx, data.ID)
		if err != nil {
			global.Logger.Error(err.Error())
			reply.Reply(errcode.ErrServer)
			ctx.Abort()
			return
		}

		if !ok {
			reply.Reply(errcodes.UserNotFound)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// GetTokenContent 从当前上下文中获取保存的 Content 内容
func GetTokenContent(ctx *gin.Context) (*model.Content, bool) {
	value, ok := ctx.Get(global.PrivateSetting.Token.AuthorizationKey)
	if !ok {
		return nil, false
	}
	return value.(*model.Content), true
}
