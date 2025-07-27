package task

import (
	"im/dao"
	"im/global"
	"im/model"
	"im/model/chat"
	"im/model/chat/server"

	"github.com/XYYSWK/Lutils/pkg/utils"
)

// CreateNotify 创建群通知
func CreateNotify(accessToken string, accountID, relationID int64, msgContent string, msgExtend *model.MsgExtend) func() {
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	members, err := dao.Database.DB.GetGroupMembers(ctx, relationID)
	if err != nil {
		global.Logger.Error(err.Error())
	}
	return func() {
		global.ChatMap.SendMany(members, chat.ServerCreateNotify, server.CreateNotify{
			EnToken:    utils.EncodeMD5(accessToken),
			AccountID:  accountID,
			RelationID: relationID,
			MsgContent: msgContent,
			MsgExtend:  msgExtend,
		})
	}
}

// UpdateNotify 更新群通知
func UpdateNotify(accessToken string, accountID, relationID int64, msgContent string, msgExtend *model.MsgExtend) func() {
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	members, err := dao.Database.DB.GetGroupMembers(ctx, relationID)
	if err != nil {
		global.Logger.Error(err.Error())
	}
	return func() {
		global.ChatMap.SendMany(members, chat.ServerUpdateNotify, server.CreateNotify{
			EnToken:    utils.EncodeMD5(accessToken),
			AccountID:  accountID,
			RelationID: relationID,
			MsgContent: msgContent,
			MsgExtend:  msgExtend,
		})
	}
}
