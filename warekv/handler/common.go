package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"ware-kv/warekv/manager"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
)

func Get(c *gin.Context) {
	_, val := findKeyAndValue(c)
	if !isKVEffective(c, val) {
		return
	}
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  val.GetValue(),
	})
}

func Delete(c *gin.Context) {
	key, val := findKeyAndValue(c)
	if !isKVEffective(c, val) {
		return
	}
	del(key)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func set(key *storage.Key, newVal storage.Value, expireTime int64) {
	if expireTime != 0 {
		newVal.WithExpireTime(expireTime)
		time.AfterFunc(time.Duration(expireTime) * time.Second, func() {
			del(key)
		})
	}
	storage.GetWareTable().Set(key, newVal)
	setNotify(key, newVal)
}

func del(key *storage.Key) {
	storage.GetWareTable().Delete(key)
	deleteNotify(key)
}

func setNotify(key *storage.Key, newVal storage.Value) {
	go manager.GetSubscribeCenter().Notify(key.GetKey(), newVal.GetValue(), manager.CallbackSetEvent)
}

func deleteNotify(key *storage.Key) {
	go manager.GetSubscribeCenter().Notify(key.GetKey(), nil, manager.CallbackDeleteEvent)
}

func findKeyAndValue(c *gin.Context) (*storage.Key, storage.Value) {
	paramKey := c.Param("key")
	if paramKey == "" {
		log.Println("key is null")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Key should not be null!",
		})
	}
	key := storage.MakeKey(paramKey)
	val := storage.GetWareTable().Get(key)
	return key, val
}

func isKVEffective(c *gin.Context, val storage.Value) bool {
	if val == nil {
		log.Println("key is not existed")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.KeyNotExisted,
		})
		return false
	}
	if !val.IsAlive() {
		log.Println("key is dead")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.KeyHasDeleted,
		})
		return false
	}
	if val.IsExpired() {
		log.Println("key has been expired")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.KeyHasExpired,
		})
		return false
	}
	return true
}
