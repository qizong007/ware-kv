package storage

import (
	"sync"
	"time"
)

const (
	defaultWriteQueueCap     = 256 // 默认写请求缓存容量
	defaultWriteTickInterval = 100 * time.Millisecond
)

type writeEvent int

const (
	SetEvent = iota
	DeleteEvent
)

// Shard 表分片
type Shard struct {
	table      map[string]Value
	rw         sync.RWMutex
	writeQueue chan *writeReq // 写请求缓存队列
	wqTicker   *time.Ticker   // 写队列的定时器
	gcTasks    chan string    // gc task 任务队列（存key）
	gcTicker   *time.Ticker   // gc task 的定时器
}

// 写请求
type writeReq struct {
	event writeEvent
	key   *Key
	value Value
}

func newShard() *Shard {
	shard := &Shard{
		table:      make(map[string]Value),
		writeQueue: make(chan *writeReq, defaultWriteQueueCap),
		gcTasks:    make(chan string, defaultGCTaskCap),
		wqTicker:   time.NewTicker(defaultWriteTickInterval),
		gcTicker:   time.NewTicker(defaultGCTickInterval),
	}
	go shard.scheduledBatchCommit()
	return shard
}

func (s *Shard) Get(key *Key) Value {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.table[key.GetKey()]
}

func (s *Shard) Set(key *Key, val Value) {
	s.writeQueue <- &writeReq{
		event: SetEvent,
		key:   key,
		value: val,
	}
}

func (s *Shard) Delete(key *Key) {
	s.writeQueue <- &writeReq{
		event: DeleteEvent,
		key:   key,
	}
}

func (s *Shard) scheduledBatchCommit() {
	for {
		select {
		case <-s.wqTicker.C: // 批量写入
			if len(s.writeQueue) == 0 {
				continue
			}
			// 批量写入
			s.rw.Lock()
			for entry := range s.writeQueue {
				key := entry.key.GetKey()
				switch entry.event {
				case SetEvent:
					s.table[key] = entry.value
				case DeleteEvent:
					s.table[key].DeleteValue()
					s.gcTasks <- key
				}
				if len(s.writeQueue) == 0 {
					break
				}
			}
			s.rw.Unlock()
		case <-s.gcTicker.C: // 批量清扫
			if len(s.gcTasks) == 0 {
				continue
			}
			s.rw.Lock()
			for key := range s.gcTasks {
				delete(s.table, key)
				if len(s.gcTasks) == 0 {
					break
				}
			}
			s.rw.Unlock()
		}
	}
}
