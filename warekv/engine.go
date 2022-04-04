package warekv

import (
	"github.com/qizong007/ware-kv/warekv/machine"
	"github.com/qizong007/ware-kv/warekv/manager"
	"github.com/qizong007/ware-kv/warekv/storage"
	"log"
)

type WareEngine struct {
	wTable          storage.KVTable
	subscribeCenter *manager.SubscribeCenter
	info            *machine.Info
	// TODO Hot Key Sampling
}

var engine *WareEngine

func Default() *WareEngine {
	return New(nil)
}

type WareEngineOption struct {
	Shard       *storage.ShardOption           `yaml:"Shard"`
	GC          *storage.WareGCOption          `yaml:"GC"`
	Subscriber  *manager.SubscribeCenterOption `yaml:"Subscriber"`
	MachineInfo *machine.WareInfoOption        `yaml:"MachineInfo"`
	Cache       *CacheOption                   `yaml:"Cache"`
}

type CacheOption struct {
	Open     bool   `yaml:"Open"`
	Strategy string `yaml:"Strategy"`
	MaxBytes uint64 `yaml:"MaxBytes"`
}

func DefaultCacheOption() *CacheOption {
	return &CacheOption{Open: false}
}

func DefaultOption() *WareEngineOption {
	return &WareEngineOption{
		Shard:       storage.DefaultShardOption(),
		GC:          storage.DefaultWareGCOption(),
		Subscriber:  manager.DefaultSubscribeCenterOption(),
		MachineInfo: machine.DefaultWareInfoOption(),
		Cache:       DefaultCacheOption(),
	}
}

func New(option *WareEngineOption) *WareEngine {
	engine = &WareEngine{}
	if option == nil {
		engine.wTable = storage.NewWareTable(nil, nil)
		engine.subscribeCenter = manager.NewSubscribeCenter(nil)
		engine.info = machine.NewWareInfo(nil)
		storage.GlobalTable = engine.wTable
	} else {
		if option.Cache != nil && option.Cache.Open {
			engine.wTable = newCache(option.Cache)
		} else {
			engine.wTable = storage.NewWareTable(option.Shard, option.GC)
		}
		storage.GlobalTable = engine.wTable
		engine.subscribeCenter = manager.NewSubscribeCenter(option.Subscriber)
		engine.info = machine.NewWareInfo(option.MachineInfo)
	}
	return engine
}

func newCache(option *CacheOption) storage.KVTable {
	switch option.Strategy {
	case storage.LRUStrategy:
		log.Println("Cache Eviction Strategy ->", storage.LRUStrategy)
		return storage.NewLRUCache(int64(option.MaxBytes))
	case storage.LFUStrategy:
		log.Println("Cache Eviction Strategy ->", storage.LFUStrategy)
		return storage.NewLFUCache(int64(option.MaxBytes))
	default:
		log.Println("Cache Eviction Strategy ->", storage.LRUStrategy)
		return storage.NewLRUCache(int64(option.MaxBytes))
	}
}

func (e *WareEngine) Close() {
	engine.info.Close()
	engine.subscribeCenter.Close()
	engine.wTable.Close()
}

func Engine() *WareEngine {
	return engine
}

func (e *WareEngine) Get(key *storage.Key) storage.Value {
	return e.wTable.Get(key)
}

func (e *WareEngine) Set(key *storage.Key, value storage.Value) {
	e.wTable.Set(key, value)
}

func (e *WareEngine) SetInTime(key *storage.Key, value storage.Value) {
	e.wTable.SetInTime(key, value)
}

func (e *WareEngine) Delete(key *storage.Key) {
	e.wTable.Delete(key)
}

func (e *WareEngine) Subscribe(option *manager.SubscribeManifest) {
	e.subscribeCenter.Subscribe(option)
}

func (e *WareEngine) Notify(key string, newVal interface{}, event int) {
	e.subscribeCenter.Notify(key, newVal, event)
}
