package router

import (
	_ "im/docs"
	"im/global"
	"im/middlewares"
	"im/routers"

	socketio "github.com/googollee/go-socket.io"

	"github.com/XYYSWK/Lutils/pkg/app"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func NewRouter() (*gin.Engine, *socketio.Server) {
	//创建一个新的路由
	r := gin.New()
	r.Use(middlewares.Cors(), middlewares.GinLogger(), middlewares.Recovery(true))
	root := r.Group("api", middlewares.LogBody(), middlewares.ParseToAuth())
	{
		root.GET("swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
		root.GET("/ping", func(ctx *gin.Context) {
			reply := app.NewResponse(ctx)
			// 使用...将切片展开为多个参数
			global.Logger.Info("ping", middlewares.ErrLogMsg(ctx)...)
			reply.Reply(nil, "pong")
		})
		rg := routers.Routers
		rg.User.Init(root)
		rg.Email.Init(root)
		rg.Account.Init(root)
		rg.Application.Init(root)
		rg.Group.Init(root)
		rg.Setting.Init(root)
		rg.File.Init(root)
		rg.Notify.Init(root)
	}
	return r, routers.Routers.Chat.Init(r)
}
