package settings

import (
	"im/global"

	"github.com/XYYSWK/Lutils/pkg/app"
)

type page struct {
}

// Init 分页器初始化
func (page) Init() {
	global.Page = app.InitPage(global.PublicSetting.Page.DefaultPageSize,
		global.PublicSetting.Page.MaxPageSize,
		global.PublicSetting.Page.PageKey,
		global.PublicSetting.Page.PageSizeKey)
}
