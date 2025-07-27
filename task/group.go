package task

import (
	"im/dao"
	"im/global"
	"im/model/chat"
	"im/model/chat/server"

	"github.com/XYYSWK/Lutils/pkg/utils"
)

func TransferGroup(accessToken string, accountID, relationID int64) func() {
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	// 获取群中所有的成员的ID
	members, err := dao.Database.Redis.GetAllAccountsByRelationID(ctx, relationID)
	if err != nil {
		global.Logger.Error(err.Error())
	}
	return func() {
		global.ChatMap.SendMany(members, chat.ServerGroupTransferred, server.TransferGroup{
			EnToken:   utils.EncodeMD5(accessToken),
			AccountID: accountID,
		})
	}
}

func DissolveGroup(accessToken string, relationID int64) func() {
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	members, err := dao.Database.Redis.GetAllAccountsByRelationID(ctx, relationID)
	if err != nil {
		global.Logger.Error(err.Error())
	}
	return func() {
		global.ChatMap.SendMany(members, chat.ServerGroupDissolved, server.DissolveGroup{
			EnToken:    utils.EncodeMD5(accessToken),
			RelationID: relationID,
		})
	}
}

func InviteGroup(accessToken string, accountID, relationID int64) func() {
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	members, err := dao.Database.Redis.GetAllAccountsByRelationID(ctx, relationID)
	if err != nil {
		global.Logger.Error(err.Error())
	}
	return func() {
		global.ChatMap.SendMany(members, chat.ServerInviteAccount, server.InviteGroup{
			EnToken:   utils.EncodeMD5(accessToken),
			AccountID: accountID,
		})
	}
}

func QuitGroup(accessToken string, accountID, relationID int64) func() {
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	members, err := dao.Database.Redis.GetAllAccountsByRelationID(ctx, relationID)
	if err != nil {
		global.Logger.Error(err.Error())
	}
	return func() {
		global.ChatMap.SendMany(members, chat.ServerQuitGroup, server.QuitGroup{
			EnToken:   utils.EncodeMD5(accessToken),
			AccountID: accountID,
		})
	}
}
