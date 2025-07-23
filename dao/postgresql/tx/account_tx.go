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

		// 建立关系（自己与自己的好友关系）
		// 步骤4：建立账户与自身的好友关系
		// 允许系统将用户与自身的关系视为特殊好友关系处理
		var relationID int64
		err = tool.DoThat(err, func() error {
			relationID, err = queries.CreateFriendRelation(ctx, &db.CreateFriendRelationParams{
				Account1ID: arg.ID,
				Account2ID: arg.ID,
			})
			return err
		})

		// 步骤5：为该关系创建默认设置
		err = tool.DoThat(err, func() error {
			return queries.CreateSetting(ctx, &db.CreateSettingParams{
				AccountID:  arg.ID,
				RelationID: relationID,
				IsSelf:     true,
			})
		})

		// 添加自己一个人的关系到 redis
		err = tool.DoThat(err, func() error {
			return rdb.AddRelationAccount(ctx, relationID, arg.ID)
		})
		return err
	})
}

// DeleteAccountWithTx 删除账户
func (store *SqlStore) DeleteAccountWithTx(ctx context.Context, rdb *operate.RDB, accountID int64) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		// 判断该账户是不是群主
		var isLeader bool
		err = tool.DoThat(err, func() error {
			isLeader, err = queries.ExistsGroupLeaderByAccountIDWithLock(ctx, accountID)
			return err
		})
		if isLeader {
			return ErrAccountGroupLeader
		}
		// 删除好友
		var friendRelationIDs []int64
		err = tool.DoThat(err, func() error {
			friendRelationIDs, err = queries.DeleteFriendRelationsByAccountID(ctx, accountID)
			return err
		})
		// 删除群
		var groupRelationIDs []int64
		err = tool.DoThat(err, func() error {
			groupRelationIDs, err = queries.DeleteSettingsByAccountID(ctx, accountID)
			return err
		})
		// 删除账户
		err = tool.DoThat(err, func() error {
			err = queries.DeleteAccount(ctx, accountID)
			return err
		})

		// 从redis 中删除对应的关系
		// 从redis 中删除该账户的好友关系
		err = tool.DoThat(err, func() error {
			return rdb.DeleteRelations(ctx, friendRelationIDs...)
		})
		// 从redis 中删除该账户所在的群聊中的账户
		err = tool.DoThat(err, func() error {
			return rdb.DeleteAccountFromRelations(ctx, accountID, groupRelationIDs...)
		})
		return err
	})
}
