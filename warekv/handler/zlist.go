package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/warekv/ds"
	"ware-kv/warekv/manager"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
)

func SetZList(c *gin.Context) {
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

	list, ok := val.([]interface{})
	if !ok {
		log.Println("v is not map slice")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "v should be like: [{\"score\": 80, \"val\": \"Sam\"}] !",
		})
		return
	}

	newList, err := interfaceSlice2SlElementSlice(list)
	if err != nil {
		log.Println("interfaceSlice2SlElementSlice fail", err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "v should be like: [{\"score\": 80, \"val\": \"Sam\"}] !",
		})
		return
	}
	newVal := ds.MakeZList(newList)
	storage.GetWareTable().Set(key, newVal)
	go manager.GetSubscribeCenter().Notify(key.GetKey(), newList, manager.CallbackSetEvent)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func interfaceSlice2SlElementSlice(list []interface{}) ([]util.SlElement,error) {
	res := make([]util.SlElement, 0, len(list))
	for i := range list {
		e, ok := list[i].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("change to SlElement fail")
		}
		res = append(res, util.SlElement{
			Score: e["score"].(float64),
			Val:   e["val"],
		})
	}
	return res, nil
}

func GetZListLen(c *gin.Context) {
	_, val := findKeyAndValue(c)
	if !isValEffective(c, val) {
		return
	}
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  ds.Value2ZList(val).GetLen(),
	})
}
