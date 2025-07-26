package request

import "mime/multipart"

type ParamPublishFile struct {
	File       *multipart.FileHeader `form:"file" binding:"required" swaggerignore:"true"`
	RelationID int64                 `form:"relation_id" binding:"required"`
	AccountID  int64                 `form:"account_id" binding:"required"`
}
