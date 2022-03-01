package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/warekv/ds"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
)

type SetStrParam struct {
	Key        string `json:"k"`
	Val        string `json:"v"`
	ExpireTime int64  `json:"expire_time" binding:"-"`
}

func SetStr(c *gin.Context) {
	param := SetStrParam{}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	key := storage.MakeKey(param.Key)
	newVal := ds.MakeString(param.Val)

	set(key, newVal, param.ExpireTime)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func GetStrLen(c *gin.Context) {
	_, val := findKeyAndValue(c)
	if !isKVEffective(c, val) {
		return
	}
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  ds.Value2String(val).GetLen(),
	})
}
