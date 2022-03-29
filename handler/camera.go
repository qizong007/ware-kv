package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/camera"
	"ware-kv/util"
	"ware-kv/warekv/storage"
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

	go camera.GetCamera().TakePhotos([]storage.Photographer{storage.GlobalTable}, param.NeedZip)
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  "Save Worker Start...",
	})
}
