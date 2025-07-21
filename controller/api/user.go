package api

import (
	"im/errcodes"
	"im/logic"
	"im/middlewares"
	"im/model"
	"im/model/request"

	"github.com/XYYSWK/Lutils/pkg/app/errcode"

	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/gin-gonic/gin"
)

type user struct {
}

// Register 注册用户
// @Tags user
// @Summary 用户注册
// @accept application/json
// @Produce application/json
// @Param data body request.ParamRegister true "用户注册信息"
// @Success 200 {object} common.State{data=reply.ParamRegister} "1001:参数有误 1003:系统错误 2004:邮箱验证码校验失败 2006:邮箱已经注册"
// @Router /api/user/register [post]
func (user) Register(ctx *gin.Context) {
	// 1.获取参数和参数校验
	reply := app.NewResponse(ctx)
	params := new(request.ParamRegister)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcodes.PasswordNotValid.WithDetails(err.Error()))
		return
	}
	// 2.业务处理
	result, err := logic.Logics.User.Register(ctx, params.Email, params.Password, params.Code)
	// 3.返回响应
	reply.Reply(err, result)
}

// Login 用户登录
// @Tags user
// @Summary 用户登录
// @accept application/json
// @Produce application/json
// @Param data body request.ParamLogin true "用户登录信息"
// @Success 200 {object} common.State{data=reply.ParamLogin} "1001:参数错误 1003:系统错误 2001:用户不存在"
// @Router /api/user/login [post]
func (user) Login(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamLogin)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcodes.PasswordNotValid.WithDetails(err.Error()))
		return
	}
	result, err := logic.Logics.User.Login(ctx, params.Email, params.Password)
	reply.Reply(err, result)
}

// UpdateUserPassword 更新用户密码
// @Tags user
// @Summary 更新用户密码
// @Security BasicAuth
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param data body request.ParamUpdateUserPassword true "验证码和新密码"
// @Success 200 {object} common.State{} "1001:参数有误 1003:系统错误 2001:用户不存在 2004:邮箱验证码校验失败 2007:身份不存在 2008:身份验证失败"
// @Router /api/user/update/pwd [put]
func (user) UpdateUserPassword(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateUserPassword)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.User.UpdateUserPassword(ctx, content.ID, params.Code, params.NewPassword)
	reply.Reply(err)
}
