package storage

import "time"

const (
	defaultGCTaskCap      = 1024 // 默认 gc 任务缓存容量
	defaultGCTickInterval = 500 * time.Millisecond
)
