package warekv

import (
	"github.com/gin-gonic/gin"
	handler2 "ware-kv/warekv/handler"
	util2 "ware-kv/warekv/util"
)

func Register(r *gin.Engine) {
	// 添加计时器
	r.Use(util2.TimeKeeping())
	// string
	r.GET("/str/:key", handler2.GetStr)
	r.POST("/str", handler2.SetStr)
	r.PUT("/str", handler2.SetStr)
	r.DELETE("/str/:key", handler2.DeleteStr)
	r.GET("/str/:key/len", handler2.GetStrLen)
	// test
	r.GET("/ping", handler2.Ping)
}
