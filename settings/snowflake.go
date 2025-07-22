package settings

import (
	"im/global"
	"time"

	"github.com/XYYSWK/Lutils/pkg/generateID/snowflake"
)

type generateID struct {
}

func (generateID) Init() {
	var err error
	global.GenerateID, err = snowflake.Init(time.Now(), global.PublicSetting.App.MachineID)
	if err != nil {
		panic(err)
	}
}
