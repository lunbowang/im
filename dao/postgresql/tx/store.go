package tx

import (
	"context"
	"fmt"
	db "im/dao/postgresql/sqlc"
	"im/dao/redis/operate"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

// TXer 定义一个接口，用于执行事务相关操作
type TXer interface {
	// CreateAccountWithTx 创建账号并建立和自己的关系
	CreateAccountWithTx(ctx context.Context, rdb *operate.RDB, maxAccountNum int32, arg *db.CreateAccountParams) error
}

// SqlStore 用于处理数据类型
type SqlStore struct {
	*db.Queries               //嵌入 *db.Queries,可以直接访问 db.Queries 中定义的字段，不需要间接访问
	Pool        *pgxpool.Pool //连接池
}

// 通过事务执行回调函数
func (store *SqlStore) execTx(ctx context.Context, fn func(queries *db.Queries) error) error {
	// 开启一个数据事务
	tx, err := store.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted, // 设置事务隔离级别为已提交读。即事务只能看到已经提交的数据，可以防止脏读和不可重复读取，但不能防止幻读
		AccessMode:     pgx.ReadWrite,     // 设置事务访问模式为读写。即事务具有读取和写入数据的权限，可以执行对数据库进行修改的操作
		DeferrableMode: pgx.Deferrable,    // 设置事务延迟模式为可延迟。即事务可以延迟到其他事务结束后才提交，以确保事务的一致性。
	})
	if err != nil {
		return err
	}
	// 使用开启的事务创建一个查询
	q := store.WithTx(tx)
	// 调用传入的回调函数执行数据库操作
	if err := fn(q); err != nil {
		// 如果回调函数失败，回溯事务
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err:%v, rb err:%v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx) //提交事务
}
