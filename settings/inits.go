package settings

type group struct {
	Config     config
	Logger     log
	Worker     worker
	Dao        database
	TokenMaker tokenMaker
	EmailMark  mark
}

var Group = new(group)

// Inits 初始化项目
func Inits() {
	Group.Config.Init()
	Group.Dao.Init()
	Group.Logger.Init()
	Group.Worker.Init()
	Group.TokenMaker.Init()
	Group.EmailMark.Init()
}
