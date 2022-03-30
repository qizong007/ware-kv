package storage

import (
	"github.com/qizong007/ware-kv/warekv/storage/ds"
	"github.com/qizong007/ware-kv/warekv/util"
)

// Value all the data structure should implement it
type Value interface {
	GetValue() interface{}
	DeleteValue()
	IsAlive() bool
	IsExpired() bool
	WithExpireTime(t int64)
	Update()
	GetType() util.DSType
	GetBase() *ds.Base
	SetBase(*ds.Base)
}

func Value2Bitmap(val Value) *ds.Bitmap {
	return val.(*ds.Bitmap)
}

func Value2BloomFilter(val Value) *ds.BloomFilter {
	return val.(*ds.BloomFilter)
}

func Value2Counter(val Value) *ds.Counter {
	return val.(*ds.Counter)
}

func Value2List(val Value) *ds.List {
	return val.(*ds.List)
}

func Value2Lock(val Value) *ds.Lock {
	return val.(*ds.Lock)
}

func Value2Object(val Value) *ds.Object {
	return val.(*ds.Object)
}

func Value2Set(val Value) *ds.Set {
	return val.(*ds.Set)
}

func Value2String(val Value) *ds.String {
	return val.(*ds.String)
}

func Value2ZList(val Value) *ds.ZList {
	return val.(*ds.ZList)
}
