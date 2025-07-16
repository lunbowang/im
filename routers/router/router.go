package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() (*gin.Engine, error) {
	//创建一个新的路由
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return r, nil
}
