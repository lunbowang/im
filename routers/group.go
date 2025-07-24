package routers

import (
	"im/controller/api"
	"im/middlewares"

	"github.com/gin-gonic/gin"
)

type group struct {
}

func (group) Init(router *gin.RouterGroup) {
	r := router.Group("group").Use(middlewares.MustAccount())
	{
		r.POST("create", api.Apis.Group.CreateGroup)
		r.POST("invite", api.Apis.Group.InviteAccount)
		r.POST("transfer", api.Apis.Group.TransferGroup)
		r.POST("dissolve", api.Apis.Group.DissolveGroup)
		r.POST("update", api.Apis.Group.UpdateGroup)
		r.GET("list", api.Apis.Group.GetGroupList)
		r.POST("quit", api.Apis.Group.QuitGroup)
		r.POST("name", api.Apis.Group.GetGroupsByName)
		r.GET("members", api.Apis.Group.GetGroupMembers)
	}
}
