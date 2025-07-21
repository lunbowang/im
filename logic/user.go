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

// Login 登录
func (user) Login(ctx *gin.Context, emailStr, pwd string) (*reply.ParamLogin, errcode.Err) {
	// 根据邮箱获取用户信息
	userInfo, myerr := getUserInfoByEmail(ctx, emailStr)
	if myerr != nil {
		return nil, myerr
	}
	// 校验密码
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

	// 将token保存在redis中
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

// UpdateUserPassword 更新用户密码
// 参数：code 验证码 newPwd 新密码
func (user) UpdateUserPassword(ctx *gin.Context, userID int64, code, newPwd string) errcode.Err {
	userInfo, myerr := getUserInfoByID(ctx, userID)
	if myerr != nil {
		global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
		return myerr
	}
	// 校验验证码
	if !global.EmailMark.CheckCode(userInfo.Email, code) {
		return errcodes.EmailCodeNotValid
	}
	// 密码加密
	hashPassword, err := password.HashPassword(newPwd)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}

	//更新密码
	if err = dao.Database.DB.UpdateUser(ctx, &db.UpdateUserParams{
		Email:    userInfo.Email,
		Password: hashPassword,
		ID:       userInfo.ID,
	}); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}

	// 清除用户token
	if err := dao.Database.Redis.DeleteAllTokenByUser(ctx, userID); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

// getUserInfoByID 通过用户 ID 获取用户信息
// 参数：userID 用户ID
// 成功：用户信息，nil
// 失败：打印日志 errcodes.UserNotFound, errcode.ErrServer
func getUserInfoByID(ctx *gin.Context, userID int64) (*db.User, errcode.Err) {
	userInfo, err := dao.Database.DB.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errcodes.UserNotFound
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return userInfo, nil
}
