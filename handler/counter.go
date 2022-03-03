package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"ware-kv/warekv/ds"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
)

type SetCounterParam struct {
	Key        string `json:"k"`
	Val        int64  `json:"v"`
	ExpireTime int64  `json:"expire_time" binding:"-"`
}

func SetCounter(c *gin.Context) {
	param := SetCounterParam{}
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
	newVal := ds.MakeCounter(param.Val)

	set(key, newVal, param.ExpireTime)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func IncrCounter(c *gin.Context) {
	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}

	counter := ds.Value2Counter(val)
	counter.Incr()

	setNotify(key, counter)
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  counter.GetValue(),
	})
}

func IncrByCounter(c *gin.Context) {
	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}

	param := c.Param("delta")
	if param == "" {
		log.Println("delta is <nil>")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "delta is <nil>!",
		})
		return
	}

	delta, err := util.Str2Int64(c.Param("delta"))
	if err != nil {
		log.Println(err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "IncrByCounter's delta Str2Int64 fail!",
		})
		return
	}

	counter := ds.Value2Counter(val)
	counter.IncrBy(delta)

	setNotify(key, counter)
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  counter.GetValue(),
	})
}

func DecrCounter(c *gin.Context) {
	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}

	counter := ds.Value2Counter(val)
	counter.Decr()

	setNotify(key, counter)
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  counter.GetValue(),
	})
}

func DecrByCounter(c *gin.Context) {
	key, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}

	param := c.Param("delta")
	if param == "" {
		log.Println("delta is <nil>")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "delta is <nil>!",
		})
		return
	}

	delta, err := util.Str2Int64(c.Param("delta"))
	if err != nil {
		log.Println(err)
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "IncrByCounter's delta Str2Int64 fail!",
		})
		return
	}

	counter := ds.Value2Counter(val)
	counter.DecrBy(delta)

	setNotify(key, counter)
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  counter.GetValue(),
	})
}
