package api

import (
	"im/errcodes"
	"im/logic"
	"im/model/request"

	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/gin-gonic/gin"
)

type user struct {
}

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
