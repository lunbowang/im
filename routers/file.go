package routers

import (
	"im/controller/api"
	"im/middlewares"

	"github.com/gin-gonic/gin"
)

type file struct {
}

func (file) Init(router *gin.RouterGroup) {
	r := router.Group("file", middlewares.MustAccount())
	{
		r.POST("publish", api.Apis.File.PublishFile)
		r.DELETE("delete", api.Apis.File.DeleteFile)
		r.GET("getFiles", api.Apis.File.GetRelationFile)
		avatarGroup := r.Group("avatar")
		{
			avatarGroup.PUT("account", api.Apis.File.UploadAccountAvatar)
		}
		r.GET("details", api.Apis.File.GetFileDetailsByID)
	}
}
