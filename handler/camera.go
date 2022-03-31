package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qizong007/ware-kv/camera"
	"github.com/qizong007/ware-kv/util"
	"github.com/qizong007/ware-kv/warekv/manager"
	"github.com/qizong007/ware-kv/warekv/storage"
	"log"
)

type CameraSaveParam struct {
	NeedZip bool `json:"zip"`
}

func CameraSave(c *gin.Context) {
	param := CameraSaveParam{}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	if !camera.GetCamera().IsActive() {
		util.MakeResponse(c, &util.WareResponse{
			Code: util.CameraNotOpen,
		})
	}

	go camera.GetCamera().TakePhotos([]storage.Photographer{storage.GlobalTable, manager.GetSubscribeCenter()}, param.NeedZip)
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  "Save Worker Start...",
	})
}
