package router

import (
	"im/global"
	"im/middlewares"

	"github.com/XYYSWK/Lutils/pkg/app"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	//创建一个新的路由
	r := gin.New()
	r.Use(middlewares.Cors(), middlewares.GinLogger(), middlewares.Recovery(true))
	root := r.Group("api", middlewares.LogBody(), middlewares.ParseToAuth())
	{
		root.GET("/ping", func(ctx *gin.Context) {
			reply := app.NewResponse(ctx)
			// 使用...将切片展开为多个参数
			global.Logger.Info("ping", middlewares.ErrLogMsg(ctx)...)
			reply.Reply(nil, "pong")
		})
	}
	return r
}
