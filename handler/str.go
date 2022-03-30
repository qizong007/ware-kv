package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qizong007/ware-kv/tracker"
	"github.com/qizong007/ware-kv/util"
	"github.com/qizong007/ware-kv/warekv/storage"
	"github.com/qizong007/ware-kv/warekv/storage/ds"
	dstype "github.com/qizong007/ware-kv/warekv/util"
	"log"
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

	cmd := tracker.NewCreateCommand(param.Key, tracker.StringStruct, param.Val, newVal.CreateTime, param.ExpireTime)
	set(key, newVal, param.ExpireTime, cmd)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func GetStrLen(c *gin.Context) {
	_, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, dstype.StringDS) {
		return
	}

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  storage.Value2String(val).GetLen(),
	})
}
