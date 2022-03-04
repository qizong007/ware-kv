package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/util"
	"ware-kv/warekv/ds"
	"ware-kv/warekv/storage"
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

	set(key, newVal, param.ExpireTime)

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
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  ds.Value2Object(val).GetFieldByKey(c.Param("field")),
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

	filed := c.Param("field")

	if filed == "" {
		log.Println("field is null")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "field is null",
		})
		return
	}

	object := ds.Value2Object(val)
	object.SetFieldByKey(filed, param.Val)
	setNotify(key, object)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}
