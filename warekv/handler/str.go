package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/warekv/ds"
	"ware-kv/warekv/manager"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
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

	newVal := ds.MakeString(val.(string))
	storage.GetWareTable().Set(key, newVal)
	go manager.GetSubscribeCenter().Notify(key.GetKey(), newVal, manager.CallbackSetEvent)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func DeleteStr(c *gin.Context) {
	key, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}
	storage.GetWareTable().Delete(key)
	//go manager.GetSubscribeCenter().Notify(key.GetKey(), nil, manager.CallbackSetEvent)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
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
