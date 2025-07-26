package logic

import (
	"database/sql"
	"im/dao"
	db "im/dao/postgresql/sqlc"
	"im/errcodes"
	"im/global"
	"im/middlewares"
	"im/model"
	"im/pkg/gtype"

	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type file struct {
}

func (file) PublishFile(ctx *gin.Context, params model.PublishFile) (model.PublishFileReply, errcode.Err) {
	fileType, myErr := gtype.GetFileType(params.File)
	if myErr != nil {
		return model.PublishFileReply{}, errcode.ErrServer
	}
	if fileType == "file" {
		if params.File.Size > global.PublicSetting.Rules.BiggestFileSize {
			return model.PublishFileReply{}, errcodes.FileTooBig
		}
	} else {
		fileType = "img"
	}
	// todo 华为云云存储

	r, err := dao.Database.DB.CreateFile(ctx, &db.CreateFileParams{
		FileName: params.File.Filename,
		FileType: db.Filetype(fileType),
		FileSize: params.File.Size,
		Key:      "",
		Url:      "",
		RelationID: sql.NullInt64{
			Int64: params.RelationID,
			Valid: true,
		}, AccountID: sql.NullInt64{
			Int64: params.AccountID,
			Valid: true,
		},
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return model.PublishFileReply{}, errcode.ErrServer
	}
	return model.PublishFileReply{
		ID:       r.ID,
		FileType: fileType,
		FileSize: r.FileSize,
		Url:      r.Url,
		CreateAt: r.CreateAt,
	}, nil
}
