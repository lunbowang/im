package task

import (
	"im/dao"
	"im/global"
	"im/model/chat"
	"im/model/chat/server"
	"im/model/reply"
	"im/pkg/rocketmq/producer"

	"github.com/XYYSWK/Lutils/pkg/utils"
)

/*
有关消息的推送任务
*/

func PublishMsg(msg reply.ParamMsgInfoWithRly) func() {
	return func() {
		// 创建有超时事件的上下文
		ctx, cancel := global.DefaultContextWithTimeout()
		defer cancel()
		accountIDs, err := dao.Database.Redis.GetAllAccountsByRelationID(ctx, msg.RelationID)
		if err != nil {
			global.Logger.Error(err.Error())
			return
		}
		for _, accountID := range accountIDs {
			// 用户如果在线，直接将消息发送过去
			if global.ChatMap.CheckIsOnConnection(accountID) {
				global.ChatMap.Send(accountID, chat.ClientSendMsg, msg)
			} else {
				// 用户处于离线状态，将消息发送至MQ中
				producer.SendMsgToMQ(accountID, msg)
			}
		}
	}
}

// ReadMsg 推送阅读消息事件
// 参数：读者 ID，消息 Map(accountID:[]msgID)，所有 msgIDs
func ReadMsg(accessToken string, readerID int64, msgMap map[int64][]int64, allMsgIDs []int64) func() {
	return func() {
		if len(msgMap) == 0 {
			return
		}
		enToken := utils.EncodeMD5(accessToken)
		// 给发送消息者推送
		for accountID, msgIDs := range msgMap {
			global.ChatMap.Send(accountID, chat.ClientReadMsg, server.ReadMsg{
				EnToken:  enToken,
				MsgIDs:   msgIDs,
				ReaderID: readerID,
			})
		}
		// 给自己的其他设备同步
		global.ChatMap.Send(readerID, chat.ClientReadMsg, server.ReadMsg{
			EnToken:  enToken,
			MsgIDs:   allMsgIDs,
			ReaderID: readerID,
		})
	}
}

// UpdateMsgState 更新消息状态
func UpdateMsgState(accessToken string, relationID, msgID int64, msgType server.MsgType, state bool) func() {
	return func() {
		ctx, cancel := global.DefaultContextWithTimeout()
		defer cancel()
		accountIDs, err := dao.Database.Redis.GetAllAccountsByRelationID(ctx, relationID)
		if err != nil {
			global.Logger.Error(err.Error())
			return
		}
		global.ChatMap.SendMany(accountIDs, chat.ServerUpdateMsgState, server.UpdateMsgState{
			EnToken: utils.EncodeMD5(accessToken),
			MsgType: msgType,
			MsgID:   msgID,
			State:   state,
		})
	}
}
