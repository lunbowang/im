package operate

import (
	"context"

	"github.com/XYYSWK/Lutils/pkg/utils"
)

const keyGroup = "KeyGroup"

// AddRelationAccount 向群聊名单中添加人员(两个人的好友关系，相当于一个特殊的群聊)
func (r *RDB) AddRelationAccount(ctx context.Context, relationID int64, accountIDs ...int64) error {
	if len(accountIDs) == 0 {
		return nil
	}
	data := make([]interface{}, len(accountIDs))
	for i, v := range accountIDs {
		data[i] = utils.IDToString(v)
	}
	return r.rdb.SAdd(ctx, utils.LinkStr(keyGroup, utils.IDToString(relationID)), data...).Err()
}

// DeleteRelations 删除指定账号的指定好友
func (r *RDB) DeleteRelations(ctx context.Context, relationIDs ...int64) error {
	if len(relationIDs) == 0 {
		return nil
	}
	// 创建一个 redis 管道，可以在单个操作中批量执行多个命令，提高了性能
	pipe := r.rdb.TxPipeline()
	for _, relationID := range relationIDs {
		// 添加删除命令到管道中
		pipe.Del(ctx, utils.LinkStr(keyGroup, utils.IDToString(relationID)))
	}
	// 执行管道中的所有命令
	_, err := pipe.Exec(ctx)
	return err
}

// DeleteAccountFromRelations 从多个群聊中删除指定账号
func (r *RDB) DeleteAccountFromRelations(ctx context.Context, accountID int64, relationIDs ...int64) error {
	if len(relationIDs) == 0 {
		return nil
	}
	pipe := r.rdb.TxPipeline()
	for _, relationID := range relationIDs {
		pipe.SRem(ctx, utils.LinkStr(keyGroup, utils.IDToString(relationID)), utils.IDToString(accountID))
	}
	_, err := pipe.Exec(ctx)
	return err
}

// DeleteRelationAccount 从一个群聊中删除多个成员
func (r *RDB) DeleteRelationAccount(ctx context.Context, relationID int64, accountIDs ...int64) error {
	if len(accountIDs) == 0 {
		return nil
	}
	data := make([]interface{}, len(accountIDs))
	for i, v := range accountIDs {
		data[i] = utils.IDToString(v)
	}
	return r.rdb.SRem(ctx, utils.LinkStr(keyGroup, utils.IDToString(relationID)), data...).Err()
}

// GetAllAccountsByRelationID 从 redis 中获取所有 Account
func (r *RDB) GetAllAccountsByRelationID(ctx context.Context, relationID int64) ([]int64, error) {
	id := utils.IDToString(relationID)
	key := utils.LinkStr(keyGroup, id)
	accountIDStr, err := r.rdb.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	accountIDs := make([]int64, 0, len(accountIDStr))
	for _, str := range accountIDStr {
		accountID := utils.StringToIDMust(str)
		accountIDs = append(accountIDs, accountID)
	}
	return accountIDs, nil
}
