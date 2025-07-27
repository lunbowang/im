package logic

import (
	"context"
	"database/sql"
	"errors"
	"im/dao"
	db "im/dao/postgresql/sqlc"
	"im/errcodes"
	"im/global"
	"im/middlewares"
	"im/model"
	"im/model/chat/server"
	"im/model/reply"
	"im/task"
	"strings"

	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

type setting struct {
}

// UpdateNickName 更改给好友的备注昵称或自己在群组中的昵称
func (setting) UpdateNickName(ctx *gin.Context, accountID, relationID int64, nickName string) errcode.Err {
	settingInfo, err := dao.Database.DB.GetSettingByID(ctx, &db.GetSettingByIDParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errcodes.RelationNotExists
	case errors.Is(err, nil):
		if settingInfo.NickName == nickName {
			return nil
		}
		if err := dao.Database.DB.UpdateSettingNickName(ctx, &db.UpdateSettingNickNameParams{
			NickName:   nickName,
			AccountID:  accountID,
			RelationID: relationID,
		}); err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		// 向自己推送更改昵称的通知
		accessToken, _ := middlewares.GetToken(ctx.Request.Header)
		global.Worker.SendTask(task.UpdateNickName(accessToken, accountID, relationID, nickName))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

// UpdateSettingPin 更新好友或群组的pin（置顶）状态
func (setting) UpdateSettingPin(ctx *gin.Context, accountID, relationID int64, isPin bool) errcode.Err {
	settingInfo, err := dao.Database.DB.GetSettingByID(ctx, &db.GetSettingByIDParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errcodes.RelationNotExists
	case errors.Is(err, nil):
		if settingInfo.IsPin == isPin {
			return nil
		}
		if err := dao.Database.DB.UpdateSettingPin(ctx, &db.UpdateSettingPinParams{
			IsPin:      isPin,
			AccountID:  accountID,
			RelationID: relationID,
		}); err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		// 向自己推送更改指定通知
		accessToken, _ := middlewares.GetToken(ctx.Request.Header)
		global.Worker.SendTask(task.UpdateSettingState(accessToken, server.SettingPin, accountID, relationID, isPin))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

// UpdateSettingDisturb 更新好友或群组的免打扰状态
func (setting) UpdateSettingDisturb(ctx *gin.Context, accountID, relationID int64, isNotDisturb bool) errcode.Err {
	settingInfo, err := dao.Database.DB.GetSettingByID(ctx, &db.GetSettingByIDParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errcodes.RelationNotExists
	case errors.Is(err, nil):
		if settingInfo.IsNotDisturb == isNotDisturb {
			return nil
		}
		err = dao.Database.DB.UpdateSettingDisturb(ctx, &db.UpdateSettingDisturbParams{
			IsNotDisturb: isNotDisturb,
			AccountID:    accountID,
			RelationID:   relationID,
		})
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		// 向自己更新免打扰的通知
		accessToken, _ := middlewares.GetToken(ctx.Request.Header)
		global.Worker.SendTask(task.UpdateSettingState(accessToken, server.SettingNotDisturb, accountID, relationID, isNotDisturb))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

// UpdateSettingShow 更新好友或群组的是否展示的状态
func (setting) UpdateSettingShow(ctx *gin.Context, accountID, relationID int64, isShow bool) errcode.Err {
	settingInfo, err := dao.Database.DB.GetSettingByID(ctx, &db.GetSettingByIDParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errcodes.RelationNotExists
	case errors.Is(err, nil):
		if settingInfo.IsShow == isShow {
			return nil
		}
		err = dao.Database.DB.UpdateSettingShow(ctx, &db.UpdateSettingShowParams{
			IsShow:     isShow,
			AccountID:  accountID,
			RelationID: relationID,
		})
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		// 向自己推送更改展示状态的通知
		accountToken, _ := middlewares.GetToken(ctx.Request.Header)
		global.Worker.SendTask(task.UpdateSettingState(accountToken, server.SettingShow, accountID, relationID, isShow))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

// GetPins 获取当前账户所有pin的好友和群组列表
func (setting) GetPins(ctx *gin.Context, accountID int64) (*reply.ParamGetPins, errcode.Err) {
	friendData, err := dao.Database.DB.GetFriendPinSettingsOrderByPinTime(ctx, accountID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	groupData, err := dao.Database.DB.GetGroupPinSettingsOrderByPinTime(ctx, accountID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}

	result := make([]*model.SettingPin, 0, len(friendData)+len(groupData))
	for i, j := 0, 0; i < len(friendData) || j < len(groupData); {
		if i < len(friendData) && (j < len(groupData) || friendData[i].PinTime.Before(groupData[j].PinTime)) {
			v := friendData[i]
			friendInfo := &model.SettingFriendInfo{
				AccountID: accountID,
				Name:      v.AccountName,
				Avatar:    v.AccountAvatar,
			}
			result = append(result, &model.SettingPin{
				SettingPinInfo: model.SettingPinInfo{
					RelationID:   v.RelationID,
					RelationType: "friend",
					NickName:     v.NickName,
					PinTime:      v.PinTime,
				},
				FriendInfo: friendInfo,
			})
			i++
		} else {
			v := groupData[j]
			groupType := strings.Split(v.GroupType.String, ",")
			groupInfo := &model.SettingGroupInfo{
				RelationID:  v.RelationID,
				Name:        groupType[0][1:], // 去掉前面的 "
				Description: groupType[1],
				Avatar:      groupType[2][:len(groupType[2])-1], // 去掉后面的 "
			}
			result = append(result, &model.SettingPin{
				SettingPinInfo: model.SettingPinInfo{
					RelationID:   v.RelationID,
					RelationType: "group",
					NickName:     v.NickName,
					PinTime:      v.PinTime,
				},
				GroupInfo: groupInfo,
			})
			j++
		}
	}
	return &reply.ParamGetPins{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

// GetShows 获取当前账户首页显示的好友和群组列表
func (setting) GetShows(ctx *gin.Context, accountID int64) (*reply.ParamGetShows, errcode.Err) {
	friendData, err := dao.Database.DB.GetFriendShowSettingsOrderByShowTime(ctx, accountID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	groupData, err := dao.Database.DB.GetGroupShowSettingsOrderByShowTime(ctx, accountID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	result := make([]*model.Setting, 0, len(friendData)+len(groupData))
	for i, j := 0, 0; i < len(friendData) || j < len(groupData); {
		if i < len(friendData) && (j >= len(groupData) || friendData[i].LastShow.After(groupData[j].LastShow)) {
			v := friendData[i]
			friendInfo := &model.SettingFriendInfo{
				AccountID: v.AccountID,
				Name:      v.AccountName,
				Avatar:    v.AccountAvatar,
			}
			result = append(result, &model.Setting{
				SettingInfo: model.SettingInfo{
					RelationID:   v.RelationID,
					RelationType: "friend",
					NickName:     v.NickName,
					IsNotDisturb: v.IsNotDisturb,
					IsShow:       v.IsShow,
					PinTime:      v.PinTime,
					LastShow:     v.LastShow,
				},
				FriendInfo: friendInfo,
			})
			i++
		} else {
			v := groupData[j]
			groupType := strings.Split(v.GroupType.String, ",")
			groupInfo := &model.SettingGroupInfo{
				RelationID:  v.RelationID,
				Name:        groupType[0][1:],
				Description: groupType[1],
				Avatar:      groupType[2][:len(groupType[2])-1],
			}
			result = append(result, &model.Setting{
				SettingInfo: model.SettingInfo{
					RelationID:   v.RelationID,
					RelationType: "group",
					NickName:     v.NickName,
					IsNotDisturb: v.IsNotDisturb,
					IsPin:        v.IsPin,
					IsShow:       v.IsShow,
					PinTime:      v.PinTime,
					LastShow:     v.LastShow,
				},
				GroupInfo: groupInfo,
			})
			j++
		}
	}
	return &reply.ParamGetShows{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

// GetFriends 获取当前账户所有的好友
func (setting) GetFriends(ctx *gin.Context, accountID int64) (*reply.ParamGetFriends, errcode.Err) {
	data, err := dao.Database.DB.GetFriendSettingsOrderByName(ctx, accountID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	result := make([]*model.SettingFriend, 0, len(data))
	for _, v := range data {
		result = append(result, &model.SettingFriend{
			SettingInfo: model.SettingInfo{
				RelationID:   v.RelationID,
				RelationType: "friend",
				NickName:     v.NickName,
				IsNotDisturb: v.IsNotDisturb,
				IsPin:        v.IsPin,
				IsShow:       v.IsShow,
				PinTime:      v.PinTime,
				LastShow:     v.LastShow,
			},
			FriendInfo: &model.SettingFriendInfo{
				AccountID: v.AccountID,
				Name:      v.AccountName,
				Avatar:    v.AccountAvatar,
			},
		})
	}
	return &reply.ParamGetFriends{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

// DeleteFriend 删除好友关系（双向删除）
func (setting) DeleteFriend(ctx *gin.Context, accountID, relationID int64) errcode.Err {
	friendInfo, myErr := getFriendRelationByID(ctx, relationID)
	if myErr != nil {
		return myErr
	}
	if friendInfo.Account1ID != accountID && friendInfo.Account2ID != accountID {
		return errcodes.AuthPermissionsInsufficient
	}
	if err := dao.Database.DB.DeleteRelationWithTx(ctx, dao.Database.Redis, relationID); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}

	// 推送删除通知
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.DeleteRelation(accessToken, accountID, relationID))
	return nil
}

func getFriendRelationByID(ctx *gin.Context, relationID int64) (*db.GetFriendRelationByIDRow, errcode.Err) {
	relationInfo, err := dao.Database.DB.GetFriendRelationByID(ctx, relationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.RelationNotExists
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return relationInfo, nil
}

// GetFriendsByName 通过姓名模糊查询好友（好友姓名或昵称）
func (setting) GetFriendsByName(ctx *gin.Context, accountID int64, name string, limit, offset int32) (*reply.ParamGetFriendsByName, errcode.Err) {
	data, err := dao.Database.DB.GetFriendSettingsByName(ctx, &db.GetFriendSettingsByNameParams{
		AccountID: accountID,
		Limit:     limit,
		Offset:    offset,
		Name:      name,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamGetFriendsByName{List: []*model.SettingFriend{}}, nil
	}
	result := make([]*model.SettingFriend, 0, len(data))
	for _, v := range data {
		result = append(result, &model.SettingFriend{
			SettingInfo: model.SettingInfo{
				RelationID:   v.RelationID,
				RelationType: string(db.RelationtypeFriend),
				NickName:     v.NickName,
				IsNotDisturb: v.IsNotDisturb,
				IsPin:        v.IsPin,
				IsShow:       v.IsShow,
				PinTime:      v.PinTime,
				LastShow:     v.LastShow,
			},
			FriendInfo: &model.SettingFriendInfo{
				AccountID: v.AccountID,
				Name:      v.AccountName,
				Avatar:    v.AccountAvatar,
			},
		})
	}
	return &reply.ParamGetFriendsByName{
		List:  result,
		Total: data[0].Total,
	}, nil
}

// ExistsSetting 是否存在 account 和 relation 关系的联系
// 参数：accountID，relationDI
// 成功：是否存在，nil
// 失败：打印错误日志 errcode.ErrServer
func ExistsSetting(ctx context.Context, accountID, relationID int64) (bool, errcode.Err) {
	ok, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	if err != nil {
		global.Logger.Error(err.Error())
		return false, errcode.ErrServer
	}
	return ok, nil
}
