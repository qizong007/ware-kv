package tracker

import (
	"encoding/json"
	"log"
	"time"
	"ware-kv/warekv"
	"ware-kv/warekv/ds"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
)

const (
	StringStruct        = "str"
	ListStruct          = "lst"
	ZListStruct         = "zl"
	SetStruct           = "set"
	ObjectStruct        = "obj"
	CounterStruct       = "cnt"
	BitmapStruct        = "bm"
	BloomStructSpecific = "bfs"
	BloomStructFuzzy    = "bff"
	LockStruct          = "lk"
)

var struct2MakeFunc = map[string]func(interface{}) storage.Value{
	StringStruct: func(str interface{}) storage.Value {
		return ds.MakeString(str.(string))
	},
	ListStruct: func(list interface{}) storage.Value {
		return ds.MakeList(list.([]interface{}))
	},
	ZListStruct: func(l interface{}) storage.Value {
		list := l.([]interface{})
		elements := make([]util.SlElement, len(list))
		for i := range list {
			e := list[i].(map[string]interface{})
			elements[i] = util.SlElement{
				Score: e["score"].(float64),
				Val:   e["val"],
			}
		}
		return ds.MakeZList(elements)
	},
	SetStruct: func(list interface{}) storage.Value {
		return ds.MakeSet(list.([]interface{}))
	},
	ObjectStruct: func(obj interface{}) storage.Value {
		return ds.MakeObject(obj.(map[string]interface{}))
	},
	CounterStruct: func(num interface{}) storage.Value {
		return ds.MakeCounter(int64(num.(float64)))
	},
	BitmapStruct: func(num interface{}) storage.Value {
		bm := ds.MakeBitmap()
		bm.SetBit(int(num.(float64)))
		return bm
	},
	BloomStructSpecific: func(param interface{}) storage.Value {
		mp := param.(map[string]interface{})
		option := ds.BloomFilterSpecificOption{
			M: uint64(mp["M"].(float64)),
			K: uint64(mp["K"].(float64)),
		}
		return ds.MakeBloomFilterSpecific(option)
	},
	BloomStructFuzzy: func(param interface{}) storage.Value {
		mp := param.(map[string]interface{})
		option := ds.BloomFilterFuzzyOption{
			N:  uint(mp["N"].(float64)),
			Fp: mp["Fp"].(float64),
		}
		return ds.MakeBloomFilterFuzzy(option)
	},
	LockStruct: func(param interface{}) storage.Value {
		return ds.MakeLock()
	},
}

type CreateCommand struct {
	Key        string      `json:"k"`
	Structure  string      `json:"s"`
	Val        interface{} `json:"v"`
	CreateTime int64       `json:"c"`
	ExpireTime int64       `json:"e"`
}

func NewCreateCommand(key string, structure string, val interface{}, createTime, expireTime int64) *CreateCommand {
	return &CreateCommand{
		Key:        key,
		Structure:  structure,
		Val:        val,
		CreateTime: createTime,
		ExpireTime: expireTime,
	}
}

func (c *CreateCommand) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		log.Println("CreateCommand Json Marshall Fail", err)
		return ""
	}
	return string(data)
}

func (c *CreateCommand) Execute() {
	key := storage.MakeKey(c.Key)
	val := struct2MakeFunc[c.Structure](c.Val)
	if c.ExpireTime != 0 {
		now := time.Now().Unix()
		// kv is now expired
		if c.CreateTime+c.ExpireTime < now {
			return
		}
		expireTime := c.ExpireTime - (now - c.CreateTime)
		val.WithExpireTime(expireTime)
		time.AfterFunc(time.Duration(expireTime)*time.Second, func() {
			warekv.Engine().Delete(key)
		})
	}
	warekv.Engine().SetInTime(key, val)
}

func (c *CreateCommand) GetOpType() string {
	return CreateOp
}
