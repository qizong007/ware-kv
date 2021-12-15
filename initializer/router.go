package initializer

import (
	"github.com/gin-gonic/gin"
	"ware-kv/handler"
	"ware-kv/util"
)

func Register(r *gin.Engine) {
	// 添加计时器
	r.Use(util.TimeKeeping())
	// string
	r.GET("/str/:key", handler.GetStr)
	r.POST("/str", handler.SetStr)
	r.PUT("/str", handler.SetStr)
	r.DELETE("/str/:key", handler.DeleteStr)
	r.GET("/str/:key/len", handler.GetStrLen)
	// test
	r.GET("/ping", handler.Ping)
}
