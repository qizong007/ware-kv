package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"ware-kv/warekv/manager"
	"ware-kv/warekv/util"
)

func SubscribeKey(c *gin.Context) {
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
	var paramPath interface{}
	var paramEvents interface{}
	var paramRetryTimes interface{}
	var ok bool
	if paramKey, ok = optionMap["key"]; !ok {
		log.Println("key is null")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Key should not be null!",
		})
		return
	}
	if paramPath, ok = optionMap["path"]; !ok {
		log.Println("path is null")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Path should not be null!",
		})
		return
	}

	option := &manager.SubscribeOption{
		Key:          paramKey.(string),
		CallbackPath: paramPath.(string),
		ExpectEvent:  nil,
		RetryTimes:   0,
	}
	if paramEvents, ok = optionMap["expect_events"]; ok {
		list := paramEvents.([]string)
		events := make([]int, 0)
		for i := range list {
			if strings.ToLower(list[i]) == "set" {
				events = append(events, manager.CallbackSetEvent)
			}
			if strings.ToLower(list[i]) == "delete" {
				events = append(events, manager.CallbackDeleteEvent)
			}
		}
		option.ExpectEvent = events
	}
	if paramRetryTimes, ok = optionMap["retry_times"]; ok {
		option.RetryTimes = paramRetryTimes.(int)
	}

	manager.GetSubscribeCenter().Subscribe(option)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}
