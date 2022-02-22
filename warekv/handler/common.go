package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
)

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
