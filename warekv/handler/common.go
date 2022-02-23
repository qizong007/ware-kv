package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/warekv/manager"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
)

func Get(c *gin.Context) {
	_, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  val.GetValue(),
	})
}

func Delete(c *gin.Context) {
	key, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}
	storage.GetWareTable().Delete(key)
	go manager.GetSubscribeCenter().Notify(key.GetKey(), nil, manager.CallbackDeleteEvent)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
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

func isValEffective(c *gin.Context, val storage.Value) bool {
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
			Msg:  "Key has been deleted...",
		})
		return false
	}
	return true
}
