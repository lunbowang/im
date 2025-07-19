package operate

import "context"

/*
redis 中邮件地址集合的 CRUD 操作。因为邮件的地址信息需要频繁访问和更新的数据，使用 Redis 可以提高性能和响应速度。
*/

const EmailKey = "EmailKey" //email set(无序集合)的键值

// AddEmails 向redis set 中添加 Emails
func (r *RDB) AddEmails(ctx context.Context, emails ...string) error {
	if len(emails) == 0 {
		return nil
	}
	data := make([]interface{}, len(emails))
	for i, email := range emails {
		data[i] = email
	}
	// 向键值为 EmailKey 的集合中添加邮箱地址集合
	return r.rdb.SAdd(ctx, EmailKey, data...).Err()
}

// ExistEmail 检查指定的 email 是否存在于 set 中
func (r *RDB) ExistEmail(ctx context.Context, email string) (bool, error) {
	return r.rdb.SIsMember(ctx, EmailKey, email).Result()
}
