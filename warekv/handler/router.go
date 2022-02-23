package handler

import (
	"github.com/gin-gonic/gin"
	"ware-kv/warekv/util"
)

func Register(r *gin.Engine) {
	// 添加计时器
	r.Use(util.TimeKeeping())
	// string
	r.GET("/str/:key", GetStr)
	r.POST("/str", SetStr)
	r.PUT("/str", SetStr)
	r.DELETE("/str/:key", DeleteStr)
	r.GET("/str/:key/len", GetStrLen)
	// subscribe
	r.POST("/subscribe", SubscribeKey)
	// test
	r.GET("/ping", Ping)
	r.GET("/err", Err500)
}
