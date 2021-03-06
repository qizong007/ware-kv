package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qizong007/ware-kv/tracker"
	"github.com/qizong007/ware-kv/util"
	"github.com/qizong007/ware-kv/warekv/storage"
	"github.com/qizong007/ware-kv/warekv/storage/ds"
	dstype "github.com/qizong007/ware-kv/warekv/util"
	"log"
	"time"
)

type SetObjectParam struct {
	Key        string                 `json:"k"`
	Val        map[string]interface{} `json:"v"`
	ExpireTime int64                  `json:"expire_time" binding:"-"`
}

type SetObjectFieldByKeyParam struct {
	Val interface{} `json:"v"`
}

func SetObject(c *gin.Context) {
	param := SetObjectParam{}
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
	newVal := ds.MakeObject(param.Val)

	cmd := tracker.NewCreateCommand(param.Key, tracker.ObjectStruct, param.Val, newVal.CreateTime, param.ExpireTime)
	set(key, newVal, param.ExpireTime, cmd)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func GetObjectFieldByKey(c *gin.Context) {
	_, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, dstype.ObjectDS) {
		return
	}

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  storage.Value2Object(val).GetFieldByKey(c.Param("field")),
	})
}

func SetObjectFieldByKey(c *gin.Context) {
	param := SetObjectFieldByKeyParam{}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, dstype.ObjectDS) {
		return
	}

	filed := c.Param("field")

	if filed == "" {
		log.Println("field is null")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "field is null",
		})
		return
	}

	object := storage.Value2Object(val)
	wal(tracker.NewModifyCommand(key.GetKey(), tracker.ObjectSetFieldByKey, time.Now().Unix(), filed, param.Val))
	object.SetFieldByKey(filed, param.Val)
	setNotify(key, object)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func DeleteObjectFieldByKey(c *gin.Context) {
	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, dstype.ObjectDS) {
		return
	}

	filed := c.Param("field")

	if filed == "" {
		log.Println("field is null")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "field is null",
		})
		return
	}

	object := storage.Value2Object(val)
	wal(tracker.NewModifyCommand(key.GetKey(), tracker.ObjectDeleteFieldByKey, time.Now().Unix(), filed))
	object.DeleteFieldByKey(filed)
	setNotify(key, object)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}
