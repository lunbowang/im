package settings

type group struct {
	Config config
	Logger log
}

var Group = new(group)

// Inits 初始化项目
func Inits() {
	Group.Config.Init()
	Group.Logger.Init()
}
