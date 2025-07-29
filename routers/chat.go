package routers

import (
	"im/controller/api"
	chat2 "im/model/chat"
	"net/http"

	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type ws struct {
}

func (ws) Init(router *gin.Engine) *socketio.Server {
	//server := socketio.NewServer(nil) // 创建一个 socketIO 服务器对象
	// 创建一个支持轮询和 WebSocket 两种传输方式的 Socket.IO 服务器
	// 配置了允许所有来源的请求（通过 CheckOrigin 函数）
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})
	{
		// 定义处理客户端事件的函数
		server.OnConnect("/", api.Apis.Chat.Handle.Onconnect)
		server.OnError("/", api.Apis.Chat.Handle.OnError)
		server.OnDisconnect("/", api.Apis.Chat.Handle.OnDisconnect)
	}
	// 配置处理聊天相关的任务
	chatHandle(server)
	// 将 server（Socket.IO 服务器）绑定到 Gin 框架的路由上，以处理客户端发起的 WebSocket 连接请求。
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	return server
}

func chatHandle(server *socketio.Server) {
	namespace := "/chat"
	server.OnEvent(namespace, chat2.ClientSendMsg, api.Apis.Chat.Message.SendMsg) //处理客户端发送的消息
	server.OnEvent(namespace, chat2.ClientReadMsg, api.Apis.Chat.Message.ReadMsg) //处理客户端已读的消息
	server.OnEvent(namespace, chat2.ClientAuth, api.Apis.Chat.Handle.Auth)        // 处理客户端认证
	server.OnEvent(namespace, chat2.ClientTest, api.Apis.Chat.Handle.Test)        //测试接口
}

// allowOriginFunc 是一个用于检查跨域请求来源的函数
// 当前实现允许所有来源的请求，实际生产环境中应根据需求修改
var allowOriginFunc = func(r *http.Request) bool {
	return true
}
