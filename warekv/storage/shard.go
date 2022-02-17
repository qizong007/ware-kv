package storage

import (
	"sync"
	"time"
)

const (
	defaultWriteQueueCap = 256 // 默认写请求缓存容量
	defaultTickInterval  = 100 * time.Millisecond
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
	ticker     *time.Ticker   // 写队列的定时器
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
		ticker:     time.NewTicker(defaultTickInterval),
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
	// 写入写事件队列
	s.writeQueue <- &writeReq{
		event: SetEvent,
		key:   key,
		value: val,
	}
}

func (s *Shard) scheduledBatchCommit() {
	for {
		select {
		case <-s.ticker.C:
			if len(s.writeQueue) == 0 {
				continue
			}
			// 批量写入
			s.rw.Lock()
			for entry := range s.writeQueue {
				switch entry.event {
				case SetEvent:
					s.table[entry.key.GetKey()] = entry.value
				case DeleteEvent:
					s.table[entry.key.GetKey()].DeleteValue()
				}
				if len(s.writeQueue) == 0 {
					break
				}
			}
			s.rw.Unlock()
		}
	}
}

func (s *Shard) Delete(key *Key) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.writeQueue <- &writeReq{
		event: DeleteEvent,
		key:   key,
	}
}
