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
	}
}
