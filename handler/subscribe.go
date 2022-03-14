package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"ware-kv/tracker"
	"ware-kv/util"
	"ware-kv/warekv"
	"ware-kv/warekv/manager"
)

type SubscribeKeyParam struct {
	Key          string    `json:"key"`
	Path         string    `json:"path"`
	Events       *[]string `json:"expect_events" binding:"-"`
	RetryTimes   *int      `json:"retry_times" binding:"-"`
	IsPersistent *bool     `json:"is_persistent" binding:"-"`
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

	manifest := &manager.SubscribeManifest{
		Key:          param.Key,
		CallbackPath: param.Path,
		ExpectEvent:  nil,
		RetryTimes:   0,
		IsPersistent: false,
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
		manifest.ExpectEvent = events
	}
	if param.RetryTimes != nil {
		manifest.RetryTimes = *param.RetryTimes
	}
	if param.IsPersistent != nil {
		manifest.IsPersistent = *param.IsPersistent
	}

	wal(tracker.NewSubCommand(manifest))
	warekv.Engine().Subscribe(manifest)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}
