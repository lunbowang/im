package operate

import (
	"im/global"

	"github.com/XYYSWK/Lutils/pkg/utils"
	"github.com/gin-gonic/gin"
)

var UserKey = "user"

func (r *RDB) SaveUserToken(ctx *gin.Context, userID int64, tokens []string) error {
	key := utils.LinkStr(UserKey, utils.IDToString(userID))
	for _, token := range tokens {
		if err := r.rdb.SAdd(ctx, key, token).Err(); err != nil {
			return err
		}
		r.rdb.Expire(ctx, key, global.PrivateSetting.Token.AccessTokenExpire)
	}
	return nil
}

// DeleteAllTokenByUser 删除所有的用户token
func (r *RDB) DeleteAllTokenByUser(ctx *gin.Context, userID int64) error {
	key := utils.LinkStr(UserKey, utils.IDToString(userID))
	if err := r.rdb.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

// CheckUserTokenValid 检查redis中是否存在该用户token
func (r *RDB) CheckUserTokenValid(ctx *gin.Context, userID int64, token string) bool {
	key := utils.LinkStr(UserKey, utils.IDToString(userID))
	ok := r.rdb.SIsMember(ctx, key, token).Val()
	return ok
}
