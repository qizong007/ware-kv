package storage

import (
	"github.com/qizong007/ware-kv/warekv/util"
	"time"
)

const (
	defaultGCTaskCap      = 1024             // gc task queue's capacity
	gcTaskCapMin          = 256              // gc task queue's capacity
	gcTaskCapMax          = 64 * 1024 * 1024 // gc task queue's capacity
	defaultGCTickInterval = 500
	gcTickIntervalMin     = 100
	gcTickIntervalMax     = 5000
)

type WareGC struct {
	shard    *Shard
	gcTasks  chan string  // gc task's queue (for storing KEY)
	gcTicker *time.Ticker // gc task's ticker
}

func NewWareGC(share *Shard, option *WareGCOption) *WareGC {
	GCTaskCap := defaultGCTaskCap
	GCTickInterval := time.Millisecond * time.Duration(defaultGCTickInterval)
	if option != nil {
		GCTaskCap = util.SetIfHitLimit(int(option.TaskCap), gcTaskCapMin, gcTaskCapMax)
		tickInterval := util.SetIfHitLimit(int(option.TickInterval), gcTickIntervalMin, gcTickIntervalMax)
		GCTickInterval = time.Millisecond * time.Duration(tickInterval)
	}
	return &WareGC{
		shard:    share,
		gcTasks:  make(chan string, GCTaskCap),
		gcTicker: time.NewTicker(GCTickInterval),
	}
}

type WareGCOption struct {
	TaskCap      uint `yaml:"TaskCap"`
	TickInterval uint `yaml:"TickInterval"`
}

func DefaultWareGCOption() *WareGCOption {
	return &WareGCOption{
		TaskCap:      defaultGCTaskCap,
		TickInterval: defaultGCTickInterval,
	}
}

func (gc *WareGC) Commit(entry string) {
	gc.gcTasks <- entry
}

func (gc *WareGC) Close() {
	gc.shard.handleGC()
	gc.gcTicker.Stop()
	close(gc.gcTasks)
}
