package logic

import (
	"errors"
	"im/dao"
	db "im/dao/postgresql/sqlc"
	"im/errcodes"
	"im/global"
	"im/middlewares"
	"im/model"
	"im/model/reply"

	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/XYYSWK/Lutils/pkg/password"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

type user struct{}

// Register 用户注册
func (user) Register(ctx *gin.Context, emailStr, pwd, code string) (*reply.ParamRegister, errcode.Err) {
	// 判断邮箱是否已经注册过了
	if err := CheckEmailNotExists(ctx, emailStr); err != nil {
		return nil, err
	}

	//校验验证码
	if !global.EmailMark.CheckCode(emailStr, code) {
		return nil, errcodes.EmailCodeNotValid
	}

	// 密码加密
	hashPassword, err := password.HashPassword(pwd)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}

	// 将user写会数据库并返回UserInfo
	userInfo, err := dao.Database.DB.CreateUser(ctx, &db.CreateUserParams{
		Email:    emailStr,
		Password: hashPassword,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}

	// 添加邮箱到 redis
	err = dao.Database.Redis.AddEmails(ctx, emailStr)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}

	// 创建token
	accessToken, accessPayload, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.AccessTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	refreshToken, _, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.RefreshTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}

	// redis保存token
	if err = dao.Database.Redis.SaveUserToken(ctx, userInfo.ID, []string{accessToken, refreshToken}); err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}

	return &reply.ParamRegister{
		ParamUserInfo: reply.ParamUserInfo{
			ID:       userInfo.ID,
			Email:    userInfo.Email,
			CreateAt: userInfo.CreateAt,
		},
		Token: reply.ParamToken{
			AccessToken:   accessToken,
			AccessPayload: accessPayload,
			RefreshToken:  refreshToken,
		},
	}, nil
}

func (user) Login(ctx *gin.Context, emailStr, pwd string) (*reply.ParamLogin, errcode.Err) {
	userInfo, myerr := getUserInfoByEmail(ctx, emailStr)
	if myerr != nil {
		return nil, myerr
	}
	if err := password.CheckPassword(pwd, userInfo.Password); err != nil {
		return nil, errcodes.PasswordNotValid
	}

	// 创建token
	accessToken, accessPayload, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.AccessTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	refreshToken, _, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.RefreshTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if err = dao.Database.Redis.SaveUserToken(ctx, userInfo.ID, []string{accessToken, refreshToken}); err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	return &reply.ParamLogin{
		ParamUserInfo: reply.ParamUserInfo{
			ID:       userInfo.ID,
			Email:    userInfo.Email,
			CreateAt: userInfo.CreateAt,
		},
		Token: reply.ParamToken{
			AccessToken:   accessToken,
			AccessPayload: accessPayload,
			RefreshToken:  refreshToken,
		},
	}, nil
}

// getUserInfoByEmail 通过邮箱获取用户信息
// 参数：emailStr 邮箱
// 成功：用户信息，nil
// 失败：打印日志 errcodes.UserNotFound, errcode.ErrServer
func getUserInfoByEmail(ctx *gin.Context, emailStr string) (*db.User, errcode.Err) {
	userInfo, err := dao.Database.DB.GetUserByEmail(ctx, emailStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errcodes.UserNotFound
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return userInfo, nil
}
