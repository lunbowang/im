package routers

import (
	"im/controller/api"
	"im/middlewares"

	"github.com/gin-gonic/gin"
)

type account struct{}

func (account) Init(router *gin.RouterGroup) {
	r := router.Group("account")
	{
		userGroup := r.Group("").Use(middlewares.MustUser())
		{
			userGroup.POST("create", api.Apis.Account.CreateAccount)
		}
	}
}
