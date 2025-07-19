package postgresql

import (
	"context"
	db "im/dao/postgresql/sqlc"
	"im/dao/postgresql/tx"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB interface {
	db.Querier
	tx.TXer
}

func Init(dataSourceName string) DB {
	//创建连接池
	pool, err := pgxpool.Connect(context.Background(), dataSourceName)
	if err != nil {
		panic(err)
	}
	return &tx.SqlStore{Queries: db.New(pool), Pool: pool}
}
