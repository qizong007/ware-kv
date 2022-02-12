package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/warekv"
	ds2 "ware-kv/warekv/ds"
	storage2 "ware-kv/warekv/storage"
	util2 "ware-kv/warekv/util"
)

func GetStr(c *gin.Context) {
	_, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}
	util2.MakeResponse(c, &util2.WareResponse{
		Code: util2.Success,
		Val:  val.GetValue(),
	})
}

func SetStr(c *gin.Context) {
	optionMap := make(map[string]interface{})
	err := c.BindJSON(&optionMap)
	if err != nil {
		log.Println("BindJSON fail")
		util2.MakeResponse(c, &util2.WareResponse{
			Code: util2.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	var paramKey interface{}
	var val interface{}
	var ok bool
	if paramKey, ok = optionMap["k"]; !ok {
		log.Println("key is null")
		util2.MakeResponse(c, &util2.WareResponse{
			Code: util2.ParamError,
			Msg:  "Key should not be null!",
		})
		return
	}
	key := storage2.MakeKey(paramKey.(string))
	if val, ok = optionMap["v"]; !ok {
		log.Println("val is null")
		util2.MakeResponse(c, &util2.WareResponse{
			Code: util2.ParamError,
			Msg:  "Val should not be null!",
		})
		return
	}
	warekv.WTable.Set(key, ds2.MakeString(val.(string)))
	util2.MakeResponse(c, &util2.WareResponse{
		Code: util2.Success,
	})
}

func DeleteStr(c *gin.Context) {
	key, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}
	warekv.WTable.Delete(key)
	util2.MakeResponse(c, &util2.WareResponse{
		Code: util2.Success,
	})
}

func GetStrLen(c *gin.Context) {
	_, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}
	util2.MakeResponse(c, &util2.WareResponse{
		Code: util2.Success,
		Val:  ds2.Value2String(val).GetLen(),
	})
}

func findKeyAndValue(c *gin.Context) (*storage2.Key, storage2.Value) {
	paramKey := c.Param("key")
	if paramKey == "" {
		log.Println("key is null")
		util2.MakeResponse(c, &util2.WareResponse{
			Code: util2.ParamError,
			Msg:  "Key should not be null!",
		})
	}
	key := storage2.MakeKey(paramKey)
	val := warekv.WTable.Get(key)
	return key, val
}

func isValEffective(c *gin.Context, val storage2.Value) bool {
	if val == nil {
		log.Println("key is not existed")
		util2.MakeResponse(c, &util2.WareResponse{
			Code: util2.KeyNotExisted,
		})
		return false
	}
	if !val.IsAlive() {
		log.Println("key is dead")
		util2.MakeResponse(c, &util2.WareResponse{
			Code: util2.KeyHasDeleted,
			Msg:  "Key has been deleted...",
		})
		return false
	}
	return true
}
