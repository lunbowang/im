package tx

import (
	"context"
	"errors"
	db "im/dao/postgresql/sqlc"
)

var (
	ErrAccountOverNum     = errors.New("账户数量超过限制")
	ErrAccountNameExists  = errors.New("账户名已存在")
	ErrAccountGroupLeader = errors.New("账户是群主")
)

// CreateAccountWithTx 检查数量、账户名之后创建账户并建立和自己的关系
func (store *SqlStore) CreateAccountWithTx(ctx context.Context, maxAccountNum int32, arg *db.CreateAccountParams) error {
	return nil
}
