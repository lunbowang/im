package settings

import (
	"im/dao"
	"im/dao/postgresql"
	"im/dao/redis"
	"im/global"
)

type database struct {
}

// Init 数据库的持久化
func (d database) Init() {
	// mysql 初始化
	dao.Database.DB = postgresql.Init(global.PrivateSetting.Postgresql.SourceName)

	// redis持久化
	dao.Database.Redis = redis.Init(
		global.PrivateSetting.Redis.Address,
		global.PrivateSetting.Redis.Password,
		global.PrivateSetting.Redis.PoolSize,
		global.PrivateSetting.Redis.DB,
	)

}
