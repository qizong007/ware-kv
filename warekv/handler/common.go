package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"ware-kv/warekv/manager"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
)

func Get(c *gin.Context) {
	_, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  val.GetValue(),
	})
}

func Delete(c *gin.Context) {
	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
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
		time.AfterFunc(time.Duration(expireTime)*time.Second, func() {
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

func keyNull(c *gin.Context) {
	paramNull(c, "Key")
}

func paramNull(c *gin.Context, param string) {
	util.MakeResponse(c, &util.WareResponse{
		Code: util.ParamError,
		Msg:  param + " should not be null!",
	})
}

func findKeyAndValue(c *gin.Context) (*storage.Key, storage.Value, error) {
	return findKeyAndValByParam(c, "key")
}

func findKeyAndValByParam(c *gin.Context, param string) (*storage.Key, storage.Value, error) {
	paramKey := c.Param(param)
	if paramKey == "" {
		log.Println(param, "is null")
		return nil, nil, fmt.Errorf("%s", util.ErrCode2Msg[util.ParamError])
	}
	key := storage.MakeKey(paramKey)
	val := storage.GetWareTable().Get(key)
	return key, val, nil
}

func isValNil(val storage.Value) bool {
	if val == nil {
		return true
	}
	return false
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
