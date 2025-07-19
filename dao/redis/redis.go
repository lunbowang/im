package redis

import (
	"context"
	"im/dao/redis/operate"

	"github.com/go-redis/redis/v8"
)

// Init redis 初始化
func Init(Addr, Password string, PoolSize, DB int) *operate.RDB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     Addr,     //host:port
		Password: Password, //密码
		PoolSize: PoolSize, //连接池
		DB:       DB,       //默认连接数据库（0-15）
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return operate.New(rdb)
}
