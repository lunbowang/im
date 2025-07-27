package task

import (
	"im/global"
	"im/model/chat"
)

func Application(accountID int64) func() {
	return func() {
		global.ChatMap.Send(accountID, chat.ServerApplication)
	}
}
