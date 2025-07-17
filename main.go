package main

import (
	"context"
	"fmt"
	"im/global"
	"im/model/common"
	"im/routers/router"
	"im/settings"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

func main() {
	//1. 初始化项目（配置加载，日志，数据库，雪花算法等）
	settings.Inits()
	// 设置Gin框架为Release（生成）模式
	if global.PublicSetting.Server.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	// 验证邮箱的合法性
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("email", common.ValidatorEmail)
	}

	// 2. 注册路由，返回路由和Socket.IO 服务器实例
	r := router.NewRouter()
	// 3.启动服务（优雅关机）
	//http.Server 内置的 Shutdown()方法支持优雅关机
	sever := http.Server{
		Addr:           global.PublicSetting.Server.HttpPort, //端口号
		Handler:        r,                                    //路由处理器
		MaxHeaderBytes: 1 << 20,                              //最大请求头大小(1MB)
		//设置合适的 MaxHeaderBytes ，值可以确保服务器能够有效地处理请求头，避免不必要地资源浪费或潜在的安全风险
	}
	global.Logger.Info("Server is started!")
	fmt.Println("AppName:", global.PublicSetting.App.Name,
		"Version:", global.PublicSetting.App.Version,
		"Address:", global.PublicSetting.Server.HttpPort,
		"RunMode:", global.PublicSetting.Server.RunMode)

	errChan := make(chan error, 1)
	defer close(errChan) //延迟关闭错误通道

	go func() {
		//开启一个goroutine启动服务
		err := sever.ListenAndServe()
		if err != nil {
			errChan <- err //将错误发送到错误通道
		}
	}()

	// 优雅退出
	// 创建一个接收信号的通道
	quit := make(chan os.Signal, 1)                      //os.Signal 标识操作系统的信号，比如终端信号、终止信号等。
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //signal.Notify 把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	select {
	case err := <-errChan:
		global.Logger.Error(err.Error())
	case <-quit:
		global.Logger.Info("Shutdown Server ...")
		// 创建一个带超时的上下文（给几秒完成还未处理完的请求）
		ctx, cancel := context.WithTimeout(context.Background(), global.PublicSetting.Server.DefaultContextTimeout)
		defer cancel() //延迟取消上下文

		//上下文超时时间内优雅关机（将未处理完的请求处理完再关闭服务），超过超时时间退出
		if err := sever.Shutdown(ctx); err != nil {
			global.Logger.Error("Server forced to shutdown, err:" + err.Error())
		}
	}

	global.Logger.Info("Server exit")

}
