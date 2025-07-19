package dao

import (
	"im/dao/postgresql"
	"im/dao/redis/operate"
)

type database struct {
	DB    postgresql.DB
	Redis *operate.RDB
}

var Database = new(database)
