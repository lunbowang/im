package routers

import (
	"im/controller/api"

	"github.com/gin-gonic/gin"
)

type email struct {
}

func (email) Init(router *gin.RouterGroup) {
	r := router.Group("email")
	{
		r.POST("send", api.Apis.Email.SendMark)
	}
}
