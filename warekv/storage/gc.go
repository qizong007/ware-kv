package storage

import (
	"time"
)

const (
	defaultGCTaskCap      = 1024 // gc task queue's capacity
	defaultGCTickInterval = 500
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
		GCTaskCap = int(option.TaskCap)
		GCTickInterval = time.Millisecond * time.Duration(option.TickInterval)
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
