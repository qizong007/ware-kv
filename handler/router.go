package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qizong007/ware-kv/authentication"
	"github.com/qizong007/ware-kv/util"
	"log"
	"time"
)

// PUT:
//	Create/Replace the resource，and appends it to the resource group
// POST:
//  Create/Append a new element that belongs to the CURRENT resource to the resource group

func Register(r *gin.Engine, needAuth bool) {
	start := time.Now()

	// add ticker
	r.Use(util.TimeKeeping())

	// cross-origin
	r.Use(util.Cors())

	// authentication
	if needAuth {
		r.Use(authentication.Authentication())
	}

	// get machine info
	r.GET("/", Info)
	r.GET("/info", Info)

	// common kv
	r.GET("/:key", Get)
	r.DELETE("/:key", Delete)

	// string
	r.PUT("/str", SetStr)
	r.GET("/str/:key/len", GetStrLen)

	// counter
	r.PUT("/counter", SetCounter)
	r.POST("/counter/:key/incr", IncrCounter)
	r.POST("/counter/:key/incrby/:delta", IncrByCounter)
	r.POST("/counter/:key/decr", DecrCounter)
	r.POST("/counter/:key/decrby/:delta", DecrByCounter)

	// object
	r.PUT("/object", SetObject)
	r.GET("/object/:key/:field", GetObjectFieldByKey)
	r.POST("/object/:key/:field", SetObjectFieldByKey)

	// list
	r.PUT("/list", SetList)
	r.POST("/list/:key/add", AddList)
	r.GET("/list/:key/len", GetListLen)
	r.GET("/list/:key/pos/:pos", GetListByPos)
	r.GET("/list/:key/start/:left", GetListStartAt)
	r.GET("/list/:key/end/:right", GetListEndAt)
	r.GET("/list/:key/between/:left/:right", GetListBetween)
	r.DELETE("/list/:key", RemoveListElement)
	r.POST("/list/:key/rpush", RPushList)
	r.POST("/list/:key/rpop", RPopList)
	r.POST("/list/:key/lpush", LPushList)
	r.POST("/list/:key/lpop", LPopList)

	// zlist
	r.PUT("/zlist", SetZList)
	r.POST("/zlist/:key/add", AddZList)
	r.GET("/zlist/:key/len", GetZListLen)
	r.GET("/zlist/:key/pos/:pos", GetZListByPos)
	r.GET("/zlist/:key/start/:left", GetZListStartAt)
	r.GET("/zlist/:key/end/:right", GetZListEndAt)
	r.GET("/zlist/:key/between/:left/:right", GetZListBetween)
	r.GET("/zlist/:key/in/:min/:max", GetZListBetweenScores)
	r.DELETE("/zlist/:key", RemoveZListByScore)

	// set
	r.PUT("/set", SetSet)
	r.POST("/set/:key/add", AddSet)
	r.GET("/set/:key/size", GetSetSize)
	r.DELETE("/set/:key", RemoveSet)
	r.GET("/set/:key/contains", ContainsSet)
	r.GET("/set/inter/:set1/:set2", InterSet)
	r.GET("/set/union/:set1/:set2", UnionSet)
	r.GET("/set/diff/:set1/:set2", DiffSet)

	// bitmap
	r.PUT("/bitmap/:key/:num", SetBitmap)
	r.GET("/bitmap/:key/len", GetBitmapLen)
	r.GET("/bitmap/:key/between/:left/:right", GetBitCount)
	r.GET("/bitmap/:key/:num", GetBitmapBit)
	r.DELETE("/bitmap/:key/:num", ClearBitmap)

	// bloom filter
	r.PUT("/bloom", SetBloomSpecific)
	r.PUT("/bloom/fool", SetBloomFuzzy)
	r.POST("/bloom/:key", AddBloom)
	r.GET("/bloom/:key/size", GetBloomSize)
	r.GET("/bloom/:key/:data", TestBloom)
	r.GET("/bloom/:key/false_rate/:n", GetBloomFalseRate)
	r.DELETE("/bloom/:key", ClearBloom)

	// lock
	r.PUT("/lock/:key", Lock)
	r.POST("/unlock/:key", Unlock)

	// subscribe
	r.POST("/subscribe", SubscribeKey)

	// camera save
	r.POST("/camera/save", CameraSave)

	// test
	r.GET("/ping", Ping)
	r.GET("/err", Err500)

	log.Printf("Router finished loading in %s...\n", time.Since(start))
}
