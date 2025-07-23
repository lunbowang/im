package logic

import (
	"errors"
	"im/dao"
	db "im/dao/postgresql/sqlc"
	"im/dao/postgresql/tx"
	"im/errcodes"
	"im/global"
	"im/middlewares"
	"im/model"
	"im/model/common"
	"im/model/reply"

	"github.com/jackc/pgx/v4"

	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type account struct{}

// CreateAccount 创建用户
func (account) CreateAccount(ctx *gin.Context, userID int64, name, avatar, gender, signature string) (*reply.ParamCreateAccount, errcode.Err) {
	arg := &db.CreateAccountParams{
		ID:        global.GenerateID.GetID(),
		UserID:    userID,
		Name:      name,
		Avatar:    avatar,
		Gender:    db.Gender(gender),
		Signature: signature,
	}

	// 创建账户以及和自己的关系
	err := dao.Database.DB.CreateAccountWithTx(ctx, dao.Database.Redis, global.PublicSetting.Rules.AccountNumMax, arg)
	switch {
	case errors.Is(err, tx.ErrAccountOverNum):
		return nil, errcodes.AccountNumExcessive
	case errors.Is(err, tx.ErrAccountNameExists):
		return nil, errcodes.AccountNameExists
	case err == nil:
	default:
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}

	// 生成账户 token
	token, payload, err := newAccountToken(model.AccountToken, arg.ID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}

	return &reply.ParamCreateAccount{
		ParamAccountInfo: reply.ParamAccountInfo{
			ID:     arg.ID,
			Name:   name,
			Avatar: avatar,
			Gender: gender,
		}, ParamGetAccountToken: reply.ParamGetAccountToken{
			AccountToken: common.Token{
				Token:    token,
				ExpireAt: payload.ExpiredAt,
			}},
	}, nil
}

// GetAccountToken 获取账户token
func (account) GetAccountToken(ctx *gin.Context, userID, accountID int64) (*reply.ParamGetAccountToken, errcode.Err) {
	// 获取账户信息
	accountInfo, myerr := getAccountInfoByID(ctx, accountID, accountID)
	if myerr != nil {
		return nil, myerr
	}
	if accountInfo.UserID != userID {
		return nil, errcodes.AuthPermissionsInsufficient
	}

	// 创建token
	token, payload, err := newAccountToken(model.AccountToken, accountID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return &reply.ParamGetAccountToken{
		AccountToken: common.Token{
			Token:    token,
			ExpireAt: payload.ExpiredAt,
		}}, nil
}

// getAccountInfoByID 通过账号 ID 获取账号信息
func getAccountInfoByID(ctx *gin.Context, accountID, selfID int64) (*db.GetAccountByIDRow, errcode.Err) {
	// 查找账户信息
	accountInfo, err := dao.Database.DB.GetAccountByID(ctx, &db.GetAccountByIDParams{
		TargetID: accountID,
		SelfID:   selfID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errcodes.AccountNotFound
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return accountInfo, nil
}

// DeleteAccount 删除账户
func (account) DeleteAccount(ctx *gin.Context, userID, accountID int64) errcode.Err {
	accountInfo, myerr := getAccountInfoByID(ctx, accountID, accountID)
	if myerr != nil {
		return myerr
	}
	if accountInfo.UserID != userID {
		return errcodes.AuthPermissionsInsufficient
	}
	err := dao.Database.DB.DeleteAccountWithTx(ctx, dao.Database.Redis, accountID)
	switch {
	case errors.Is(err, tx.ErrAccountGroupLeader):
		return errcodes.AccountGroupLeader
	case errors.Is(err, nil):
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

// GetAccountByUserID 根据用户ID获取账户信息
func (account) GetAccountByUserID(ctx *gin.Context, userID int64) (reply.ParamGetAccountsByUserID, errcode.Err) {
	accountInfos, err := dao.Database.DB.GetAccountsByUserID(ctx, userID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return reply.ParamGetAccountsByUserID{}, errcode.ErrServer
	}
	result := make([]reply.ParamAccountInfo, len(accountInfos))
	for i, info := range accountInfos {
		result[i] = reply.ParamAccountInfo{
			ID:     info.ID,
			Name:   info.Name,
			Avatar: info.Avatar,
			Gender: string(info.Gender),
		}
	}
	return reply.ParamGetAccountsByUserID{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

// UpdateAccount 更新账户
func (account) UpdateAccount(ctx *gin.Context, accountID int64, name, gender, signature string) errcode.Err {
	err := dao.Database.DB.UpdateAccount(ctx, &db.UpdateAccountParams{
		Name:      name,
		Gender:    db.Gender(gender),
		Signature: signature,
		ID:        accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}

	// 获取token
	//accessToken,_:middlewares.GetToken(ctx.Request.Header)
	// todo 推送更新消息
	//global.Worker.SendTask(task.UpdateAccount(accessToken, accountID, name, gender, signature))
	return nil
}

// GetAccountByName 根据昵称模糊查找账号
func (account) GetAccountByName(ctx *gin.Context, accountID int64, name string, limit, offset int32) (reply.ParamGetAccountsByName, errcode.Err) {
	// 获取账户列表
	accounts, err := dao.Database.DB.GetAccountsByName(ctx, &db.GetAccountsByNameParams{
		Limit:     limit,
		Offset:    offset,
		Name:      name,
		AccountID: accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return reply.ParamGetAccountsByName{}, errcode.ErrServer
	}

	// 将数据库查询结果转换为响应数据结构
	result := make([]*reply.ParamFriendInfo, len(accounts))
	for i, info := range accounts {
		result[i] = &reply.ParamFriendInfo{
			ParamAccountInfo: reply.ParamAccountInfo{
				ID:     info.ID,
				Name:   info.Name,
				Avatar: info.Avatar,
				Gender: string(info.Gender),
			},
			RelationID: info.RelationID.Int64,
		}
	}
	return reply.ParamGetAccountsByName{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

// GetAccountByID 更具账户ID查询账户
func (account) GetAccountByID(ctx *gin.Context, accountID, selfID int64) (*reply.ParamGetAccountByID, errcode.Err) {
	info, myerr := getAccountInfoByID(ctx, accountID, selfID)
	if myerr != nil {
		return nil, myerr
	}
	return &reply.ParamGetAccountByID{
		Info: reply.ParamAccountInfo{
			ID:     info.ID,
			Name:   info.Name,
			Avatar: info.Avatar,
			Gender: string(info.Gender),
		},
		Signature:  info.Signature,
		CreateAt:   info.CreateAt,
		RelationID: info.RelationID.Int64,
	}, nil
}
