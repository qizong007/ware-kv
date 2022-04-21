package tracker

import (
	"encoding/json"
	"github.com/qizong007/ware-kv/anticorrosive"
	"github.com/qizong007/ware-kv/warekv"
	"github.com/qizong007/ware-kv/warekv/storage"
	"github.com/qizong007/ware-kv/warekv/storage/ds"
	"github.com/qizong007/ware-kv/warekv/util"
	"log"
)

const (
	CounterIncr            = "cnt_i"
	CounterDecr            = "cnt_d"
	CounterIncrBy          = "cnt_ib"
	CounterDecrBy          = "cnt_db"
	ObjectSetFieldByKey    = "obj_sf"
	ObjectDeleteFieldByKey = "obj_df"
	ListAdd                = "lst_a"
	ListSetAt              = "lst_sa"
	ListRemoveAt           = "lst_ra"
	ListRemoveVal          = "lst_rv"
	ListRPush              = "lst_rps"
	ListLPush              = "lst_lps"
	ListRPop               = "lst_rpp"
	ListLPop               = "lst_lpp"
	ZListAdd               = "zl_a"
	ZListRemoveScore       = "zl_rs"
	ZListRemoveScores      = "zl_rss"
	ZListRemoveInScore     = "zl_ris"
	SetAdd                 = "set_a"
	SetRemove              = "set_r"
	BitmapSet              = "bm_s"
	BitmapClear            = "bm_c"
	BloomFilterAdd         = "bf_a"
	BloomFilterClear       = "bf_c"
)

var str2ModifyFunc = map[string]func(storage.Value, []interface{}){
	CounterIncr: func(val storage.Value, null []interface{}) {
		val.(*ds.Counter).Incr()
	},
	CounterDecr: func(val storage.Value, null []interface{}) {
		val.(*ds.Counter).Decr()
	},
	CounterIncrBy: func(val storage.Value, params []interface{}) {
		delta := int64(params[0].(float64))
		val.(*ds.Counter).IncrBy(delta)
	},
	CounterDecrBy: func(val storage.Value, params []interface{}) {
		delta := int64(params[0].(float64))
		val.(*ds.Counter).DecrBy(delta)
	},
	ObjectSetFieldByKey: func(val storage.Value, params []interface{}) {
		filed := params[0].(string)
		val.(*ds.Object).SetFieldByKey(filed, params[1])
	},
	ObjectDeleteFieldByKey: func(val storage.Value, params []interface{}) {
		filed := params[0].(string)
		val.(*ds.Object).DeleteFieldByKey(filed)
	},
	ListAdd: func(val storage.Value, params []interface{}) {
		val.(*ds.List).Append(params[0].([]interface{}))
	},
	ListSetAt: func(val storage.Value, params []interface{}) {
		idx := int(params[0].(float64))
		_ = val.(*ds.List).SetAt(idx, params[1])
	},
	ListRemoveVal: func(val storage.Value, params []interface{}) {
		val.(*ds.List).RemoveVal(params[0])
	},
	ListRemoveAt: func(val storage.Value, params []interface{}) {
		pos := int(params[0].(float64))
		_ = val.(*ds.List).RemoveAt(pos)
	},
	ListRPush: func(val storage.Value, params []interface{}) {
		val.(*ds.List).RPush(params[0])
	},
	ListLPush: func(val storage.Value, params []interface{}) {
		val.(*ds.List).LPush(params[0])
	},
	ListRPop: func(val storage.Value, params []interface{}) {
		val.(*ds.List).RPop()
	},
	ListLPop: func(val storage.Value, params []interface{}) {
		val.(*ds.List).LPop()
	},
	ZListAdd: func(val storage.Value, params []interface{}) {
		list := params[0].([]interface{})
		elements := make([]util.SlElement, len(list))
		for i := range list {
			e := list[i].(map[string]interface{})
			elements[i] = util.SlElement{
				Score: e["score"].(float64),
				Val:   e["val"],
			}
		}
		val.(*ds.ZList).Add(elements)
	},
	ZListRemoveScore: func(val storage.Value, params []interface{}) {
		val.(*ds.ZList).RemoveScore(params[0].(float64))
	},
	ZListRemoveScores: func(val storage.Value, params []interface{}) {
		val.(*ds.ZList).RemoveScores(params[0].([]float64))
	},
	ZListRemoveInScore: func(val storage.Value, params []interface{}) {
		min := params[0].(float64)
		max := params[1].(float64)
		_ = val.(*ds.ZList).RemoveInScore(min, max)
	},
	SetAdd: func(val storage.Value, params []interface{}) {
		val.(*ds.Set).Add(params[0])
	},
	SetRemove: func(val storage.Value, params []interface{}) {
		val.(*ds.Set).Remove(params[0])
	},
	BitmapSet: func(val storage.Value, params []interface{}) {
		val.(*ds.Bitmap).SetBit(int(params[0].(float64)))
	},
	BitmapClear: func(val storage.Value, params []interface{}) {
		val.(*ds.Bitmap).ClearBit(int(params[0].(float64)))
	},
	BloomFilterAdd: func(val storage.Value, params []interface{}) {
		val.(*ds.BloomFilter).Add(params[0].(string))
	},
	BloomFilterClear: func(val storage.Value, params []interface{}) {
		val.(*ds.BloomFilter).ClearAll()
	},
}

type ModifyCommand struct {
	Key        string        `json:"k"`
	ModifyFunc string        `json:"mf"`
	UpdateTime int64         `json:"u"`
	Params     []interface{} `json:"p"`
}

func NewModifyCommand(key string, modifyFunc string, updateTime int64, params ...interface{}) *ModifyCommand {
	return &ModifyCommand{
		Key:        key,
		ModifyFunc: modifyFunc,
		UpdateTime: updateTime,
		Params:     params,
	}
}

func (c *ModifyCommand) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		log.Println("ModifyCommand Json Marshall Fail", err)
		return ""
	}
	return string(data)
}

func (c *ModifyCommand) Execute() {
	key := storage.MakeKey(c.Key)
	val := warekv.Engine().Get(key)
	isEffective, _ := anticorrosive.IsKVEffective(val)
	if !isEffective {
		log.Printf("%s is not effective now...\n", key.GetKey())
		return
	}
	str2ModifyFunc[c.ModifyFunc](val, c.Params)
}

func (c *ModifyCommand) GetOpType() string {
	return ModifyOp
}
