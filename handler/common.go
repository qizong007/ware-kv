package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qizong007/ware-kv/anticorrosive"
	"github.com/qizong007/ware-kv/tracker"
	"github.com/qizong007/ware-kv/util"
	"github.com/qizong007/ware-kv/warekv"
	"github.com/qizong007/ware-kv/warekv/storage"
	dstype "github.com/qizong007/ware-kv/warekv/util"
	"log"
	"time"
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

	wal(tracker.NewDeleteCommand(key.GetKey()))
	anticorrosive.Del(key)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func set(key *storage.Key, newVal storage.Value, expireTime int64, cmd tracker.Command) {
	wal(cmd)
	if expireTime != 0 {
		newVal.WithExpireTime(expireTime)
		time.AfterFunc(time.Duration(expireTime)*time.Second, func() {
			anticorrosive.Del(key)
		})
	}
	anticorrosive.Set(key, newVal)
}

func setInTime(key *storage.Key, newVal storage.Value, expireTime int64, cmd tracker.Command) {
	wal(cmd)
	if expireTime != 0 {
		newVal.WithExpireTime(expireTime)
		time.AfterFunc(time.Duration(expireTime)*time.Second, func() {
			anticorrosive.Del(key)
		})
	}
	anticorrosive.SetInTime(key, newVal)
}

func setNotify(key *storage.Key, newVal storage.Value) {
	anticorrosive.SetNotify(key, newVal)
}

func wal(cmd tracker.Command) {
	tracker.GetTracker().Write(cmd)
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
	val := warekv.Engine().Get(key)
	return key, val, nil
}

func isValNil(val storage.Value) bool {
	if val == nil {
		return true
	}
	return false
}

func isKVEffective(c *gin.Context, val storage.Value) bool {
	isEffective, code := anticorrosive.IsKVEffective(val)
	if !isEffective {
		util.MakeResponse(c, &util.WareResponse{
			Code: code,
		})
	}
	return isEffective
}

func isKVTypeCorrect(c *gin.Context, val storage.Value, tp dstype.DSType) bool {
	if !anticorrosive.IsKVTypeCorrect(val, tp) {
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ValueTypeError,
		})
		return false
	}
	return true
}
