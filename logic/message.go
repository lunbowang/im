package logic

import (
	"context"
	"database/sql"
	"errors"
	"im/dao"
	db "im/dao/postgresql/sqlc"
	"im/errcodes"
	"im/global"

	"github.com/XYYSWK/Lutils/pkg/app/errcode"
)

// GetMsgInfoByID 获取消息详情
// 参数：msgID 消息ID
// 成功：消息详情，nil
// 失败：打印错误日志 errcode.ErrServer,errcodes.MsgNotExists
func GetMsgInfoByID(ctx context.Context, msgID int64) (*db.GetMessageByIDRow, errcode.Err) {
	result, err := dao.Database.DB.GetMessageByID(ctx, msgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.MsgNotExists
		}
		global.Logger.Error(err.Error())
		return nil, errcode.ErrServer
	}
	return result, nil
}
