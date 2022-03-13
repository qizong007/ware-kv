package tracker

import (
	"encoding/json"
	"log"
	"time"
	"ware-kv/anticorrosive"
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
)

var struct2MakeFunc = map[string]func(interface{}) storage.Value{
	StringStruct: func(str interface{}) storage.Value {
		return ds.MakeString(str.(string))
	},
	ListStruct: func(list interface{}) storage.Value {
		return ds.MakeList(list.([]interface{}))
	},
	ZListStruct: func(list interface{}) storage.Value {
		return ds.MakeZList(list.([]util.SlElement))
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
	BitmapStruct: func(null interface{}) storage.Value {
		return ds.MakeBitmap()
	},
	BloomStructSpecific: func(param interface{}) storage.Value {
		option := param.(ds.BloomFilterSpecificOption)
		return ds.MakeBloomFilterSpecific(option)
	},
	BloomStructFuzzy: func(param interface{}) storage.Value {
		option := param.(ds.BloomFilterFuzzyOption)
		return ds.MakeBloomFilterFuzzy(option)
	},
}

type CreateCommand struct {
	CommandBase
	Structure  string      `json:"s"`
	Val        interface{} `json:"v"`
	CreateTime int64       `json:"c"`
	ExpireTime int64       `json:"e"`
}

func NewCreateCommand(key string, structure string, val interface{}, createTime, expireTime int64) *CreateCommand {
	return &CreateCommand{
		CommandBase: CommandBase{key},
		Structure:   structure,
		Val:         val,
		CreateTime:  createTime,
		ExpireTime:  expireTime,
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
			anticorrosive.Del(key)
		})
	}
	anticorrosive.Set(key, val)
}

func (c *CreateCommand) GetOpType() string {
	return CreateOp
}
