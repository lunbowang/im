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
