package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"im/global"
	"im/model/reply"

	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// SendMsgToMQ 通过 RocketMQ 发送消息，实现了一个生产者向指定的 RocketMQ 消息队列发送消息的逻辑
func SendMsgToMQ(mID int64, msg reply.ParamMsgInfoWithRly) {
	// 创建一个 RocketMQ 生产者，连接到指定的 RocketMQ 服务器
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{fmt.Sprintf("%s:%d", global.PrivateSetting.RocketMQ.Addr, global.PrivateSetting.RocketMQ.Port)}))
	if err != nil {
		panic(fmt.Sprintf("生成 Producer 失败：%s", err))
	}
	// 启动生产者，如果生成者启动失败，则抛出异常并停止程序
	if err := p.Start(); err != nil {
		panic(err)
	}
	// 构建消息的唯一标识（UID），用于识别和定位该消息的接收者
	uID := fmt.Sprintf("accountID_%d", mID)
	// 消息内容
	sendMSg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("序列化消息失败", err)
		return
	}

	// 使用 SendSync 方法发送消息，uID 是消息的主题(topic)，SystemMsg 是消息的内容
	// primitive.NewMessage 将消息打包成 RocketMQ 消息
	res, err := p.SendSync(context.Background(), primitive.NewMessage(uID, sendMSg))
	if err != nil {
		fmt.Println("发送消息", err)
	} else {
		fmt.Println("发送成功，res：", res.String())
	}
	// 关闭生产者，如果关闭时发生错误，则抛出异常
	if err := p.Shutdown(); err != nil {
		panic(err)
	}
}
