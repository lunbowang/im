package middlewares

import (
	"bytes"
	"im/global"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/*
自定义的 logger 中间件，不使用默认的 logger 中间件
*/

const Body = "body"

// ErrLogMsg 日志数据
func ErrLogMsg(ctx *gin.Context) []zap.Field {
	var body string
	data, ok := ctx.Get(Body)
	if ok {
		body = string(data.([]byte))
	}
	path := ctx.Request.URL.Path
	query := ctx.Request.URL.RawQuery
	fields := []zap.Field{
		zap.Int("status", ctx.Writer.Status()),            //记录响应的状态码
		zap.String("method", ctx.Request.Method),          //记录请求方法
		zap.String("path", path),                          //记录请求的路径
		zap.String("query", query),                        //记录请求的原始查询参数
		zap.String("ip", ctx.ClientIP()),                  //记录客户端的IP地址
		zap.String("user_agent", ctx.Request.UserAgent()), //记录客户端的user_agent
		zap.String("body", body),                          //记录请求的主体数据
	}
	return fields
}

// LogBody 读取 body 内容缓存下来，为之后打印日志做准备（读取请求体得内容并将其存储在Gin上下文中）
// LogBody 是一个Gin框架的中间件，用于读取HTTP请求体并将其缓存到上下文
// 由于HTTP请求体是一次性读取的流，该中间件允许在后续处理中多次访问请求体内容
// 常见用途包括日志记录、请求验证和调试等场景
func LogBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 读取请求体的全部内容到字节切片
		bodyBytes, _ := io.ReadAll(ctx.Request.Body)
		// 关闭原始请求体（ReadAll会自动关闭，但显式关闭更清晰），以确保资源的正确释放
		_ = ctx.Request.Body.Close()
		// 重新创建一个可读取的请求体
		// 使用bytes.NewBuffer创建一个新的字节缓冲区，并通过io.NopCloser包装为ReadCloser接口
		// 创建一个新地可读取的请求主体，并将之前读取的 bodyBytes 作为内容，最后将其设置回 ctx.Request.Body。
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		// 将请求体内容存储到Gin上下文中，键名为"body"
		// 后续中间件或处理函数可以通过ctx.Get("body")获取请求体内容
		// 使用一个新的缓冲区来存储请求主体的内容，而不是直接读取原始的请求主体。这样我们就可以在不影响原始请求主体的情况下，对请求主体的内容进行处理和修改
		ctx.Set("body", bodyBytes)
		// 继续处理后续中间件和请求处理函数
		ctx.Next()
	}
}

// GinLogger 接收gin 框架默认的日志，在处理每个请求时记录相关的请求信息到日志中去
func GinLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//记录请求开始的时间，用于计算请求处理耗时
		start := time.Now()

		//获取请求的路径和查询参数
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		//将控制权交给后面的中间件或处理程序执行
		//当后续处理完成后，会继续执行当前函数的剩余代码
		ctx.Next()

		//计算请求处理所耗费的时间
		cost := time.Since(start)
		global.Logger.Info(path,
			zap.Int("status", ctx.Writer.Status()),                                //HTTP响应码
			zap.String("method", ctx.Request.Method),                              //HTTP请求方法
			zap.String("path", path),                                              //请求路径
			zap.String("query", query),                                            //请求查询参数
			zap.String("ip", ctx.ClientIP()),                                      //客户端IP地址
			zap.String("user-agent", ctx.Request.UserAgent()),                     //客户端User-Agent
			zap.String("error", ctx.Errors.ByType(gin.ErrorTypePrivate).String()), //私有错误信息
			zap.Duration("cost", cost),                                            //请求耗时
		)
	}
}
