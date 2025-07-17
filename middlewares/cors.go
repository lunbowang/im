package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method               // GET \ POST \ PUT \ DELETE ...
		origin := ctx.Request.Header.Get("Origin") // 获取请求的 Origin 头部的值，Origin 头部在跨域请求中很重要，会告诉服务器请求的来源域
		if origin != "" {
			// 接受客户端发送的 Origin(重要)
			ctx.Header("Access-Control-Allow-Origin", origin)
			// 服务器支持的所有跨域请求的方法
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, HEAD, PUT")
			// 允许跨域设置可以返回其他字段，可以自定义字段
			ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,content-type,Authorization,Content-Length,X-CSRF-AccessToken,AccessToken,session, token")
			// 允许浏览器（客户端）可以解析的头部（重要）
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, token")
			// 允许客户端传递校验信息比如 cookie（重要）
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}

		// 处理跨域请求时的预检请求，确认服务端是否支持跨域请求，并返回相应的响应
		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, "options ok")
			return
		}
		ctx.Next()
	}
}
