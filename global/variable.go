package global

import (
	"im/manager"
	"im/model/config"
	"im/pkg/emailMark"

	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/XYYSWK/Lutils/pkg/generateID/snowflake"
	"github.com/XYYSWK/Lutils/pkg/goroutine/work"
	"github.com/XYYSWK/Lutils/pkg/logger"
	"github.com/XYYSWK/Lutils/pkg/token"
	upload "github.com/XYYSWK/Lutils/pkg/upload/obs"
)

var (
	PublicSetting  config.PublicConfig  //Public 配置
	PrivateSetting config.PrivateConfig //Private 配置
	Page           *app.Page            //分页
	Logger         *logger.Log          //日志
	Worker         *work.Worker         // 工作池
	TokenMaker     token.MakerToken     // token
	EmailMark      *emailMark.EmailMark // 验证码
	GenerateID     *snowflake.Snowflake //snowflake 雪花算法生成的 ID
	ChatMap        *manager.ChatMap     // 聊天链接管理器
	OBS            upload.OBS
)
