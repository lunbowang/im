package api

import (
	"im/logic"
	"im/model/request"

	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type email struct {
}

// SendMark 发送验证码(邮件)
// @Tags email
// @Summary 发送验证码(邮件)
// @accept application/json
// @Produce application/json
// @Param data body request.ParamSendEmail true "email"
// @Success 200 {object} common.State{} "1001:参数有误 1003:系统错误 2006:邮箱已经注册 2003:邮件发送过于频繁，请稍后再试"
// @Router /api/email/send [post]

func (email) SendMark(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := &request.ParamSendEmail{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	err := logic.Logics.Email.SendMark(params.Email)
	reply.Reply(err)
}
