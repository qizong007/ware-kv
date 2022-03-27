package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"ware-kv/tracker"
	"ware-kv/util"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/storage/ds"
	dstype "ware-kv/warekv/util"
)

type SetSetParam struct {
	Key        string        `json:"k"`
	Val        []interface{} `json:"v"`
	ExpireTime int64         `json:"expire_time" binding:"-"`
}

type CommonSetParam struct {
	Element interface{} `json:"e"`
}

func SetSet(c *gin.Context) {
	param := SetSetParam{}
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
	newVal := ds.MakeSet(param.Val)

	cmd := tracker.NewCreateCommand(param.Key, tracker.SetStruct, param.Val, newVal.CreateTime, param.ExpireTime)
	set(key, newVal, param.ExpireTime, cmd)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func GetSetSize(c *gin.Context) {
	_, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, dstype.SetDS) {
		return
	}

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  storage.Value2Set(val).GetSize(),
	})
}

func AddSet(c *gin.Context) {
	param := CommonSetParam{}
	err := c.BindJSON(&param)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	if param.Element == nil {
		log.Println("AddSetParam is <nil>!")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param is <nil>!",
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
	if !isKVTypeCorrect(c, val, dstype.SetDS) {
		return
	}

	st := storage.Value2Set(val)

	wal(tracker.NewModifyCommand(key.GetKey(), tracker.SetAdd, time.Now().Unix(), param.Element))
	st.Add(param.Element)
	setNotify(key, st)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func RemoveSet(c *gin.Context) {
	param := CommonSetParam{}
	err := c.BindJSON(&param)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	if param.Element == nil {
		log.Println("RemoveSetParam is <nil>!")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param is <nil>!",
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
	if !isKVTypeCorrect(c, val, dstype.SetDS) {
		return
	}

	st := storage.Value2Set(val)

	wal(tracker.NewModifyCommand(key.GetKey(), tracker.SetRemove, time.Now().Unix(), param.Element))
	st.Remove(param.Element)
	setNotify(key, st)

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
	})
}

func ContainsSet(c *gin.Context) {
	param := CommonSetParam{}
	err := c.BindJSON(&param)
	if err != nil {
		log.Println("BindJSON fail")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param bind json fail!",
		})
		return
	}

	if param.Element == nil {
		log.Println("ContainsSetParam is <nil>!")
		util.MakeResponse(c, &util.WareResponse{
			Code: util.ParamError,
			Msg:  "Param is <nil>!",
		})
		return
	}

	_, val, err := findKeyAndValue(c)
	if err != nil {
		keyNull(c)
		return
	}
	if !isKVEffective(c, val) {
		return
	}
	if !isKVTypeCorrect(c, val, dstype.SetDS) {
		return
	}

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  storage.Value2Set(val).Contains(param.Element),
	})
}

func InterSet(c *gin.Context) {
	set1, set2, ok := get2SetIsOk(c)
	if !ok {
		return
	}

	inter := set1.Intersect(set2).GetValue()

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  inter,
	})
}

func UnionSet(c *gin.Context) {
	set1, set2, ok := get2SetIsOk(c)
	if !ok {
		return
	}

	union := set1.Union(set2).GetValue()

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  union,
	})
}

func DiffSet(c *gin.Context) {
	set1, set2, ok := get2SetIsOk(c)
	if !ok {
		return
	}

	diff := set1.Diff(set2).GetValue()

	util.MakeResponse(c, &util.WareResponse{
		Code: util.Success,
		Val:  diff,
	})
}

func get2SetIsOk(c *gin.Context) (*ds.Set, *ds.Set, bool) {
	_, val1, err := findKeyAndValByParam(c, "set1")
	if err != nil {
		keyNull(c)
		return nil, nil, false
	}
	if !isKVEffective(c, val1) {
		return nil, nil, false
	}
	if !isKVTypeCorrect(c, val1, dstype.SetDS) {
		return nil, nil, false
	}

	_, val2, err := findKeyAndValByParam(c, "set2")
	if err != nil {
		keyNull(c)
		return nil, nil, false
	}
	if !isKVEffective(c, val2) {
		return nil, nil, false
	}
	if !isKVTypeCorrect(c, val2, dstype.SetDS) {
		return nil, nil, false
	}

	set1 := storage.Value2Set(val1)
	set2 := storage.Value2Set(val2)
	return set1, set2, true
}
