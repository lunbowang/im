package task

import (
	"im/dao"
	"im/global"
)

/*
有关 account 的 任务
*/

// UpdateEmail 更新邮箱时，通知用户(user)的每个账户(account)
func UpdateEmail(accessToken string, userID int64, email string) func() {
	return func() {
		ctx, cancel := global.DefaultContextWithTimeout()
		defer cancel()
		_, err := dao.Database.DB.GetAcountIDsByUserID(ctx, userID)
		if err != nil {
			global.Logger.Error(err.Error())
			return
		}
	}
}
