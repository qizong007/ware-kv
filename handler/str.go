package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/ds"
	"ware-kv/global"
	"ware-kv/storage"
	"ware-kv/util"
)

func GetStr(c *gin.Context) {
	_, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  val.GetValue(),
	})
}

func SetStr(c *gin.Context) {
	optionMap := make(map[string]interface{})
	err := c.BindJSON(&optionMap)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	var paramKey interface{}
	var val interface{}
	var ok bool
	if paramKey, ok = optionMap["k"]; !ok {
		log.Println("key is null")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Key should not be null!",
		})
		return
	}
	key := storage.MakeKey(paramKey.(string))
	if val, ok = optionMap["v"]; !ok {
		log.Println("val is null")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Val should not be null!",
		})
		return
	}
	global.WTable.Set(key, ds.MakeString(val.(string)))
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func DeleteStr(c *gin.Context) {
	_, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  val.GetValue(),
	})
}

func GetStrLen(c *gin.Context) {
	_, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  ds.Value2String(val).GetLen(),
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
	val := global.WTable.Get(key)
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
