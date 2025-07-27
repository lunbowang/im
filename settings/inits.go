package settings

type group struct {
	Config     config
	Logger     log
	Page       page
	Worker     worker
	Dao        database
	GenerateID generateID
	TokenMaker tokenMaker
	EmailMark  mark
	Chat       chat
	OBS        obs
	Load       load
}

var Group = new(group)

// Inits 初始化项目
func Inits() {
	Group.Config.Init()
	Group.Logger.Init()
	Group.Page.Init()
	Group.Worker.Init()
	Group.Dao.Init()
	Group.GenerateID.Init()
	Group.TokenMaker.Init()
	Group.EmailMark.Init()
	Group.Chat.Init()
	Group.OBS.Init()
	Group.Load.Init()
}
