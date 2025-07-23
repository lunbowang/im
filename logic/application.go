package logic

import (
	"errors"
	"im/dao"
	db "im/dao/postgresql/sqlc"
	"im/errcodes"
	"im/global"
	"im/middlewares"
	"im/model/reply"

	"github.com/jackc/pgx/v4"

	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type application struct {
}

// CreateApplication 发起好友申请
func (a application) CreateApplication(ctx *gin.Context, accountID1, accountID2 int64, msg string) errcode.Err {
	// 判断两个 accountID 是否一样，不能自己给自己发送好友申请
	if accountID1 == accountID2 {
		return errcodes.ApplicationNotValid
	}

	// 判断是否已经是好友了
	id1, id2 := sortID(accountID1, accountID2)
	exist, err := dao.Database.DB.ExistsFriendRelation(ctx, &db.ExistsFriendRelationParams{
		Account1ID: id1,
		Account2ID: id2,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if exist {
		return errcodes.RelationExists
	}

	// 创建申请
	err = dao.Database.DB.CreateApplicationTx(ctx, &db.CreateApplicationParams{
		Account1ID: accountID1,
		Account2ID: accountID2,
		ApplyMsg:   msg,
	})

	switch {
	case errors.Is(err, errcodes.ApplicationExists):
		return errcodes.ApplicationExists
	case errors.Is(err, nil):
		// todo 提示对方有新的申请消息
		//global.Worker.SendTask(task.A)
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

// DeleteApplication 删除好友申请
func (application) DeleteApplication(ctx *gin.Context, accountID1, accountID2 int64) errcode.Err {
	apply, myerr := getApplication(ctx, accountID1, accountID2)
	if myerr != nil {
		return errcode.ErrServer
	}
	if apply.Account1ID != accountID1 {
		return errcodes.AuthPermissionsInsufficient
	}
	if err := dao.Database.DB.DeleteApplication(ctx, &db.DeleteApplicationParams{
		Account1ID: accountID1,
		Account2ID: accountID2,
	}); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	return nil
}

// 获取好友请求信息
func getApplication(ctx *gin.Context, accountID1, accountID2 int64) (*db.Application, errcode.Err) {
	apply, err := dao.Database.DB.GetApplicationByID(ctx, &db.GetApplicationByIDParams{
		Account1ID: accountID1,
		Account2ID: accountID2,
	})
	switch {
	case errors.Is(err, nil):
		return apply, nil
	case errors.Is(err, pgx.ErrNoRows):
		return nil, errcodes.ApplicationNotExists
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
}

// RefuseApplication 被申请者拒绝好友申请
func (application) RefuseApplication(ctx *gin.Context, accountID1, accountID2 int64, refuseMsg string) errcode.Err {
	apply, myerr := getApplication(ctx, accountID1, accountID2)
	if myerr != nil {
		return myerr
	}
	if apply.Status == db.ApplicationstatusValue2 {
		return errcodes.ApplicationRepeatOpt
	}
	if err := dao.Database.DB.UpdateApplication(ctx, &db.UpdateApplicationParams{
		RefuseMsg:  refuseMsg,
		Status:     db.ApplicationstatusValue2,
		Account1ID: accountID1,
		Account2ID: accountID2,
	}); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	return nil
}

// AcceptApplication 同意好友申请
func (application) AcceptApplication(ctx *gin.Context, accountID1, accountID2 int64) errcode.Err {
	apply, myerr := getApplication(ctx, accountID1, accountID2)
	if myerr != nil {
		return myerr
	}
	if apply.Status == db.ApplicationstatusValue1 {
		return errcodes.ApplicationRepeatOpt
	}
	accountInfo1, myerr := getAccountInfoByID(ctx, accountID1, accountID1)
	if myerr != nil {
		return myerr
	}
	accountInfo2, myerr := getAccountInfoByID(ctx, accountID2, accountID2)
	if myerr != nil {
		return myerr
	}

	_, err := dao.Database.DB.AcceptApplicationTx(ctx, dao.Database.Redis, accountInfo1, accountInfo2)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}

	// todo 推送消息
	return nil
}

// ListApplications 获取用户的所有申请信息
func (application) ListApplications(ctx *gin.Context, accountID int64, limit, offset int32) (reply.ParamListApplications, errcode.Err) {
	list, err := dao.Database.DB.GetApplications(ctx, &db.GetApplicationsParams{
		Limit:     limit,
		Offset:    offset,
		AccountID: accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return reply.ParamListApplications{}, errcode.ErrServer
	}
	if len(list) == 0 {
		return reply.ParamListApplications{List: make([]*reply.ParamApplicationInfo, 0)}, nil
	}
	data := make([]*reply.ParamApplicationInfo, len(list))
	for i, row := range list {
		name, avatar := row.Account1Name, row.Account1Avatar
		if row.Account1ID == accountID {
			name, avatar = row.Account2Name, row.Account2Avatar
		}
		data[i] = &reply.ParamApplicationInfo{
			AccountID1: row.Account1ID,
			AccountID2: row.Account2ID,
			ApplyMsg:   row.ApplyMsg,
			Refuse:     row.RefuseMsg,
			Status:     string(row.Status),
			CreateAt:   row.CreateAt,
			UpdateAt:   row.UpdateAt,
			Name:       name,
			Avatar:     avatar,
		}
	}
	return reply.ParamListApplications{
		List:  data,
		Total: list[0].Total,
	}, nil
}
