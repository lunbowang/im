package tx

import (
	"context"
	"errors"
	db "im/dao/postgresql/sqlc"
	"im/dao/redis/operate"
	"im/pkg/tool"
)

var (
	ErrAccountOverNum     = errors.New("账户数量超过限制")
	ErrAccountNameExists  = errors.New("账户名已存在")
	ErrAccountGroupLeader = errors.New("账户是群主")
)

// CreateAccountWithTx 检查数量、账户名之后创建账户并建立和自己的关系
func (store *SqlStore) CreateAccountWithTx(ctx context.Context, rdb *operate.RDB, maxAccountNum int32, arg *db.CreateAccountParams) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		var accountNum int32

		// 步骤1：检查用户已有的账户数量
		err = tool.DoThat(err, func() error {
			accountNum, err = queries.CountAccountsByUserID(ctx, arg.UserID)
			return err
		})
		if accountNum >= maxAccountNum {
			// 账户数量超出限制
			return ErrAccountOverNum
		}

		// 步骤2：检查账户名是否已存在
		var exists bool
		err = tool.DoThat(err, func() error {
			exists, err = queries.ExistsAccountByNameAndUserID(ctx, &db.ExistsAccountByNameAndUserIDParams{
				UserID: arg.UserID,
				Name:   arg.Name,
			})
			return err
		})

		if exists {
			// 账户名已被使用
			return ErrAccountNameExists
		}

		// 步骤3：创建新账户
		err = tool.DoThat(err, func() error {
			return queries.CreateAccount(ctx, arg)
		})

		// todo 建立关系（自己与自己的好友关系）
		// 步骤4：建立账户与自身的好友关系
		// 允许系统将用户与自身的关系视为特殊好友关系处理
		//var relationID int64
		//err = tool.DoThat(err, func() error {
		//	relationID, err = queries.CreateFriendRelation(ctx, &db.CreateFriendRelationParams{
		//		Account1ID: arg.ID,
		//		Account2ID: arg.ID,
		//	})
		//	return err
		//})

		// todo 步骤5：为该关系创建默认设置
		//err = tool.DoThat(err, func() error {
		//	return queries.CreateSetting(ctx, &db.CreateSettingParams{
		//		AccountID:  arg.ID,
		//		RelationID: relationID,
		//		IsSelf:     true,
		//	})
		//})

		// todo 添加自己一个人的关系到 redis
		//err=tool.DoThat(err, func() error {
		//	return rdb.Add
		//})
		return err
	})
}
