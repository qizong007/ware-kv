package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"ware-kv/util"
	"ware-kv/warekv/manager"
)

type SubscribeKeyParam struct {
	Key        string    `json:"key"`
	Path       string    `json:"path"`
	Events     *[]string `json:"expect_events" binding:"-"`
	RetryTimes *int      `json:"retry_times" binding:"-"`
}

func SubscribeKey(c *gin.Context) {
	param := SubscribeKeyParam{}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	option := &manager.SubscribeOption{
		Key:          param.Key,
		CallbackPath: param.Path,
		ExpectEvent:  nil,
		RetryTimes:   0,
	}
	if param.Events != nil {
		list := *param.Events
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
	if param.RetryTimes != nil {
		option.RetryTimes = *param.RetryTimes
	}

	manager.GetSubscribeCenter().Subscribe(option)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}
