package storage

import (
	"time"
)

const (
	defaultGCTaskCap      = 1024 // 默认 gc 任务缓存容量
	defaultGCTickInterval = 500 * time.Millisecond
)

type WareGC struct {
	shard    *Shard
	gcTasks  chan string  // gc task 任务队列（存key）
	gcTicker *time.Ticker // gc task 的定时器
}

func NewWareGC(share *Shard, option *WareGCOption) *WareGC {
	GCTaskCap := defaultGCTaskCap
	GCTickInterval := defaultGCTickInterval
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

func (gc *WareGC) Commit(entry string) {
	gc.gcTasks <- entry
}

func (gc *WareGC) Close() {
	gc.shard.handleGC()
	gc.gcTicker.Stop()
	close(gc.gcTasks)
}
