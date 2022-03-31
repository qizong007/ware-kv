package camera

import (
	"encoding/json"
	"github.com/qizong007/ware-kv/warekv/manager"
	"github.com/qizong007/ware-kv/warekv/storage"
	"github.com/qizong007/ware-kv/warekv/storage/ds"
	dstype "github.com/qizong007/ware-kv/warekv/util"
)

const (
	magicSwitchLen = 1
	createTimeLen  = 8
	flagLen        = 1
	keysNumLen     = 4
)

type MetaInfo struct {
	IsZip      bool
	CreateTime int64
}

func reduceMetaInfo(data []byte) *MetaInfo {
	start := magicHeadLen + wareKVVersionLen
	meta := data[start : start+metaDataLen]
	magicSwitch := meta[0]
	createTime := dstype.BytesToInt64(meta[magicSwitchLen : magicSwitchLen+createTimeLen])
	return &MetaInfo{
		IsZip:      (magicSwitch & zipFlag) > 0,
		CreateTime: createTime,
	}
}

func reduceContent(data []byte) {
	n := len(data)
	cur := 0
	for {
		if cur >= n {
			return
		}
		// table flag (1 byte)
		flag := data[cur]
		// keys num (4 bytes)
		cur += flagLen
		keysNumBytes := data[cur : cur+keysNumLen]
		keysNum := dstype.BytesToInt(keysNumBytes)
		cur += keysNumLen
		switch flag {
		case storage.TableFlag:
			num := reduceKVTableView(data[cur:], keysNum)
			cur += num
		case storage.SubscribeCenterFlag:
			num := reduceSubscribeCenterView(data[cur:], keysNum)
			cur += num
		default:
			return
		}
	}
}

func reduceKVTableView(data []byte, keyNum int) int {
	if len(data) == 0 {
		return 0
	}
	cur := 0
	for i := 0; i < keyNum; i++ {
		// resolve for type (1 byte)
		tipe := data[cur]
		cur++
		// resolve for key len (4 byte)
		keyLen := dstype.BytesToInt(data[cur : cur+4])
		cur += 4
		// resolve for key (keyLen byte)
		key := string(data[cur : cur+keyLen])
		cur += keyLen
		// resolve for base len (4 byte)
		baseLen := dstype.BytesToInt(data[cur : cur+4])
		cur += 4
		// resolve for base json (value len byte)
		baseJson := string(data[cur : cur+baseLen])
		cur += baseLen
		// resolve for value len (4 byte)
		valueLen := dstype.BytesToInt(data[cur : cur+4])
		cur += 4
		// resolve for value json (value len byte)
		valueJson := string(data[cur : cur+valueLen])
		cur += valueLen

		resolveKVPair(tipe, key, baseJson, valueJson)
	}
	return cur
}

var type2Value = map[uint8]func(string) storage.Value{
	dstype.StringDS: func(valueJson string) storage.Value {
		str := ""
		_ = json.Unmarshal([]byte(valueJson), &str)
		value := ds.MakeString(str)
		return value
	},
	dstype.CounterDS: func(valueJson string) storage.Value {
		var cnt int64
		_ = json.Unmarshal([]byte(valueJson), &cnt)
		value := ds.MakeCounter(cnt)
		return value
	},
	dstype.ObjectDS: func(valueJson string) storage.Value {
		var obj map[string]interface{}
		_ = json.Unmarshal([]byte(valueJson), &obj)
		value := ds.MakeObject(obj)
		return value
	},
	dstype.ListDS: func(valueJson string) storage.Value {
		var list []interface{}
		_ = json.Unmarshal([]byte(valueJson), &list)
		value := ds.MakeList(list)
		return value
	},
	dstype.ZListDS: func(valueJson string) storage.Value {
		var list []dstype.SlElement
		_ = json.Unmarshal([]byte(valueJson), &list)
		value := ds.MakeZList(list)
		return value
	},
	dstype.SetDS: func(valueJson string) storage.Value {
		var list []interface{}
		_ = json.Unmarshal([]byte(valueJson), &list)
		value := ds.MakeSet(list)
		return value
	},
	dstype.BitmapDS: func(valueJson string) storage.Value {
		var list []uint64
		_ = json.Unmarshal([]byte(valueJson), &list)
		value := ds.MakeBitmapFromList(list)
		return value
	},
	dstype.BloomFilterDS: func(valueJson string) storage.Value {
		var bfView dstype.BloomView
		_ = json.Unmarshal([]byte(valueJson), &bfView)
		value := ds.MakeBloomFilterFromView(&bfView)
		return value
	},
	dstype.LockDS: func(valueJson string) storage.Value {
		value := ds.MakeLock()
		return value
	},
}

func resolveKVPair(tipe uint8, key string, baseJson string, valueJson string) {
	value := type2Value[tipe](valueJson)
	var base ds.Base
	_ = json.Unmarshal([]byte(baseJson), &base)
	value.SetBase(&base)
	storage.GlobalTable.SetInTime(storage.MakeKey(key), value)
}

func reduceSubscribeCenterView(data []byte, keyNum int) int {
	if len(data) == 0 {
		return 0
	}
	cur := 0
	for i := 0; i < keyNum; i++ {
		// resolve for key len (4 byte)
		keyLen := dstype.BytesToInt(data[cur : cur+4])
		cur += 4
		// resolve for key (keyLen byte)
		key := string(data[cur : cur+keyLen])
		cur += keyLen
		// resolve for value len (4 byte)
		valueLen := dstype.BytesToInt(data[cur : cur+4])
		cur += 4
		// resolve for value json (value len byte)
		valueJson := string(data[cur : cur+valueLen])
		cur += valueLen

		resolveSubKVPair(key, valueJson)
	}
	return cur
}

func resolveSubKVPair(key, valueJson string) {
	var callbackPlanOption manager.CallbackPlanOption
	_ = json.Unmarshal([]byte(valueJson), &callbackPlanOption)
	manifest := manager.SubscribeManifest{
		Key:          key,
		CallbackPath: callbackPlanOption.CallbackPath,
		ExpectEvent:  callbackPlanOption.Events,
		RetryTimes:   callbackPlanOption.RetryTimes,
		IsPersistent: callbackPlanOption.IsPersistent,
	}
	manager.GetSubscribeCenter().Subscribe(&manifest)
}
