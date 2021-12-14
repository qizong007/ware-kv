package initializer

import (
	"github.com/gin-gonic/gin"
	"ware-kv/handler"
)

func Register(r *gin.Engine) {
	// string
	r.GET("/str/:key", handler.GetStr)
	// test
	r.GET("/ping", handler.Ping)
}
