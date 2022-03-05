package warekv

import (
	"ware-kv/warekv/machine"
	"ware-kv/warekv/manager"
	"ware-kv/warekv/storage"
)

type WareEngine struct {
	wTable          *storage.WareTable
	subscribeCenter *manager.SubscribeCenter
	info            *machine.Info
	// closer
	// camera(RDB)
	// 热点采样
}

var engine *WareEngine

func Default() *WareEngine {
	return New(nil)
}

type WareEngineOption struct {
	Shard       storage.ShardOption           `yaml:"Shard"`
	GC          storage.WareGCOption          `yaml:"GC"`
	Subscriber  manager.SubscribeCenterOption `yaml:"Subscriber"`
	MachineInfo machine.WareInfoOption        `yaml:"MachineInfo"`
}

func New(option *WareEngineOption) *WareEngine {
	engine = &WareEngine{}
	if option == nil {
		engine.wTable = storage.NewWareTable(nil, nil)
		engine.subscribeCenter = manager.NewSubscribeCenter(nil)
		engine.info = machine.NewWareInfo(nil)
	} else {
		engine.wTable = storage.NewWareTable(&option.Shard, &option.GC)
		engine.subscribeCenter = manager.NewSubscribeCenter(&option.Subscriber)
		engine.info = machine.NewWareInfo(&option.MachineInfo)
	}
	engine.start()
	return engine
}

func (e *WareEngine) start() {
	engine.wTable.Start()
	engine.subscribeCenter.Start()
	engine.info.Start()
}

func (e *WareEngine) Close() {
	engine.info.Close()
	engine.subscribeCenter.Close()
	engine.wTable.Close()
}

func Get() *WareEngine {
	return engine
}
