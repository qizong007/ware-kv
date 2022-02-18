package storage

import "time"

const (
	defaultGCTaskCap      = 1024 // 默认 gc 任务缓存容量
	defaultGCTickInterval = 500 * time.Millisecond
)

type WareGC struct {
	gcTasks  chan string  // gc task 任务队列（存key）
	gcTicker *time.Ticker // gc task 的定时器
}

func NewWareGC() *WareGC {
	return &WareGC{
		gcTasks:  make(chan string, defaultGCTaskCap),
		gcTicker: time.NewTicker(defaultGCTickInterval),
	}
}

func (gc *WareGC) Commit(entry string)  {
	gc.gcTasks <- entry
}
