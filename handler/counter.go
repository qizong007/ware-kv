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

	cmd := tracker.NewCreateCommand(param.Key, tracker.CounterStruct, param.Val, newVal.CreateTime, param.ExpireTime)
	set(key, newVal, param.ExpireTime, cmd)

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
	if !isKVTypeCorrect(c, val, dstype.CounterDS) {
		return
	}

	counter := storage.Value2Counter(val)
	wal(tracker.NewModifyCommand(key.GetKey(), tracker.CounterIncr, time.Now().Unix()))
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
	if !isKVTypeCorrect(c, val, dstype.CounterDS) {
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

	counter := storage.Value2Counter(val)
	wal(tracker.NewModifyCommand(key.GetKey(), tracker.CounterIncrBy, time.Now().Unix(), delta))
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
	if !isKVTypeCorrect(c, val, dstype.CounterDS) {
		return
	}

	counter := storage.Value2Counter(val)
	wal(tracker.NewModifyCommand(key.GetKey(), tracker.CounterDecr, time.Now().Unix()))
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
	if !isKVTypeCorrect(c, val, dstype.CounterDS) {
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

	counter := storage.Value2Counter(val)
	wal(tracker.NewModifyCommand(key.GetKey(), tracker.CounterDecrBy, time.Now().Unix(), delta))
	counter.DecrBy(delta)

	setNotify(key, counter)
	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  counter.GetValue(),
	})
}
