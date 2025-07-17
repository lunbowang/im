package middlewares

import (
	"fmt"
	"im/global"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/XYYSWK/Lutils/pkg/email"
	"github.com/XYYSWK/Lutils/pkg/times"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/*
自定义的 recovery 异常捕获中间件
*/

// Recovery 捕获项目可能出现的 panic，并向开发者发送异常信息的邮件
func Recovery(stack bool) gin.HandlerFunc {
	// 初始化邮箱发送器，配置SMTP连接信息
	defaultMailer := email.NewEmail(&email.SMTPInfo{
		Port:     global.PrivateSetting.Email.Port,
		IsSSL:    global.PrivateSetting.Email.IsSSL,
		Host:     global.PrivateSetting.Email.Host,
		UserName: global.PrivateSetting.Email.Username,
		Password: global.PrivateSetting.Email.Password,
		From:     global.PrivateSetting.Email.From,
	})
	return func(ctx *gin.Context) {
		//使用defer-recover 机制捕获panic
		defer func() {
			if err := recover(); err != nil {
				// 检查连接是否断开，因为这并不是真正需要进行恐慌堆栈跟踪的情况
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connections reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 将请求对象转换为字节切片,用于日志记录和邮箱通知
				httpRequest, _ := httputil.DumpRequest(ctx.Request, false)

				//尝试获取请求体（如果之前通过中间件保存过）
				var body string
				data, ok := ctx.Get(Body)
				if ok {
					body = string(data.([]byte))
				}

				//发送异常通知邮件给开发者
				sendErr := defaultMailer.SendMail( // 短信通知
					global.PrivateSetting.Email.To,
					fmt.Sprintf("异常抛出，发生时间：%v\n", time.Now().Format(times.LayoutDate)),
					fmt.Sprintf("错误信息：%s\n请求信息：%s\n请求body:%s\n调用堆栈信息：%s\n", err, string(httpRequest), body, string(debug.Stack())),
				)

				//记录邮件发送错误
				if sendErr != nil {
					global.Logger.Error(fmt.Sprintf("email.SendMail Error: %v", sendErr.Error()))
				}

				// 处理客户端连接断开的特殊情况
				if brokenPipe {
					global.Logger.Error(ctx.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)))

					// 如果连接已断开，就无法写入状态
					ctx.Error(err.(error)) //将错误信息与上下文关联
					ctx.Abort()            //阻止调用后续的处理函数
					return
				}

				//根据配置决定是否记录堆栈信息
				if stack {
					global.Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack()))) // 记录当前 goroutine 的堆栈跟踪信息到日志中
				} else { //不需要记录到堆栈信息
					global.Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("body", body))
				}

				// 返回500内部服务器错误给客户端
				ctx.AbortWithStatus(http.StatusInternalServerError) //阻止调用后续的处理函数，并返回“服务器内部错误”的状态码
			}
		}()
		ctx.Next()
	}
}
