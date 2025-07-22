package api

import (
	"im/errcodes"
	"im/global"
	"im/logic"
	"im/middlewares"
	"im/model"
	"im/model/request"

	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type account struct{}

// CreateAccount 创建账号
// @Tags account
// @Summary 创建账号
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param data body request.ParamCreateAccount true "创建账号信息"
// @Success 200 {object} common.State{data=reply.ParamCreateAccount} "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2012:账户名已经存在"
// @Router /api/account/create [post]

func (account) CreateAccount(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamCreateAccount)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok && content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Account.CreateAccount(ctx, content.ID, params.Name, global.PublicSetting.Rules.DefaultAvatarURL, params.Gender, params.Signature)
	reply.Reply(err, result)
}
