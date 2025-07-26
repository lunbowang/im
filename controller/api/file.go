package api

import (
	"im/global"
	"im/logic"
	"im/middlewares"
	"im/model"
	"im/model/request"
	"im/pkg/gtype"

	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type file struct {
}

// PublishFile 上传文件
// @Tags file
// @Summary 上传文件
// @accept multipart/form-data
// @Param Authorization header string true "Bearer 账户令牌"
// @Param file formData request.ParamPublishFile true "文件"
// @Success 200 {object} common.State{data=reply.ParamPublishFile} "1001:参数有误 1003:系统错误 8001:存储失败"
// @Router /api/file/publish [post]
func (file) PublishFile(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamPublishFile)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	fileType, myErr := gtype.GetFileType(params.File)
	if myErr != nil {
		global.Logger.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		reply.Reply(errcode.ErrServer)
		return
	}
	if fileType != "img" && fileType != "png" && fileType != "jpg" {
		fileType = "file"
	}
	result, err := logic.Logics.File.PublishFile(ctx, model.PublishFile{
		File:       params.File,
		RelationID: params.RelationID,
		AccountID:  params.AccountID,
	})
	reply.Reply(err, result)
}
