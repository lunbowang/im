package settings

import (
	"im/global"
	"im/manager"
)

type chat struct {
}

func (chat) Init() {
	global.ChatMap = manager.NewChatMap()
}
