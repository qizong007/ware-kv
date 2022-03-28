package handler

import (
	"github.com/gin-gonic/gin"
	"ware-kv/camera"
	"ware-kv/util"
	"ware-kv/warekv/storage"
)

func CameraSave(c *gin.Context) {
	go camera.GetCamera().TakePhotos([]storage.Photographer{storage.GlobalTable})
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  "Save Worker Start...",
	})
}

func CameraDevelop(c *gin.Context) {
	camera.GetCamera().DevelopPhotos()
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}
