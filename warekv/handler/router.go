package handler

import (
	"github.com/gin-gonic/gin"
	"ware-kv/warekv/util"
)

// PUT:  替换/创建指定的资源，并将其追加到相应的资源组中
// POST: 把指定的资源当做一个资源组，并在其下创建/追加一个新的元素，使其隶属于当前资源

func Register(r *gin.Engine) {
	// 添加计时器
	r.Use(util.TimeKeeping())

	// 获取系统信息
	r.GET("/", Info)
	r.GET("/info", Info)

	// common kv
	r.GET("/:key", Get)
	r.DELETE("/:key", Delete)

	// string
	r.POST("/str", SetStr)
	r.PUT("/str", SetStr)
	r.GET("/str/:key/len", GetStrLen)

	// object
	r.POST("/object", SetObject)
	r.PUT("/object", SetObject)
	r.GET("/object/:key/:field", GetObjectFieldByKey)
	r.POST("/object/:key/:field", SetObjectFieldByKey)

	// list
	r.POST("/list", SetList)
	r.PUT("/list", SetList)
	r.POST("/list/:key/add", AddList)
	r.GET("/list/:key/len", GetListLen)
	r.GET("/list/:key/pos/:pos", GetListByPos)
	r.GET("/list/:key/start/:left", GetListStartAt)
	r.GET("/list/:key/end/:right", GetListEndAt)
	r.GET("/list/:key/between/:left/:right", GetListBetween)
	r.DELETE("/list/:key", RemoveListElement)

	// zlist
	r.POST("/zlist", SetZList)
	r.PUT("/zlist", SetZList)
	r.POST("/zlist/:key/add", AddZList)
	r.GET("/zlist/:key/len", GetZListLen)
	r.GET("/zlist/:key/pos/:pos", GetZListByPos)
	r.GET("/zlist/:key/start/:left", GetZListStartAt)
	r.GET("/zlist/:key/end/:right", GetZListEndAt)
	r.GET("/zlist/:key/between/:left/:right", GetZListBetween)
	r.DELETE("/zlist/:key", RemoveZListByScore)

	// set
	r.POST("/set", SetSet)
	r.PUT("/set", SetSet)
	r.POST("/set/:key/add", AddSet)
	r.GET("/set/:key/size", GetSetSize)
	r.DELETE("/set/:key", RemoveSet)
	r.GET("/set/:key/contains", ContainsSet)
	r.GET("/set/inter/:set1/:set2", InterSet)
	r.GET("/set/union/:set1/:set2", UnionSet)
	r.GET("/set/diff/:set1/:set2", DiffSet)

	// subscribe
	r.POST("/subscribe", SubscribeKey)

	// test
	r.GET("/ping", Ping)
	r.GET("/err", Err500)
}
