package warekv

import (
	"ware-kv/warekv/machine"
	"ware-kv/warekv/manager"
	"ware-kv/warekv/storage"
)

type WareEngine struct {
	wTable          storage.KVTable
	subscribeCenter *manager.SubscribeCenter
	info            *machine.Info
	// TODO camera(RDB)
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
}

func DefaultOption() *WareEngineOption {
	return &WareEngineOption{
		Shard:       storage.DefaultShardOption(),
		GC:          storage.DefaultWareGCOption(),
		Subscriber:  manager.DefaultSubscribeCenterOption(),
		MachineInfo: machine.DefaultWareInfoOption(),
	}
}

func New(option *WareEngineOption) *WareEngine {
	engine = &WareEngine{}
	if option == nil {
		engine.wTable = storage.NewWareTable(nil, nil)
		engine.subscribeCenter = manager.NewSubscribeCenter(nil)
		engine.info = machine.NewWareInfo(nil)
	} else {
		engine.wTable = storage.NewWareTable(option.Shard, option.GC)
		engine.subscribeCenter = manager.NewSubscribeCenter(option.Subscriber)
		engine.info = machine.NewWareInfo(option.MachineInfo)
	}
	return engine
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
