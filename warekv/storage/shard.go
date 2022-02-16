package storage

import (
	"fmt"
	"sync"
)

const (
	DefaultWriteQueueCap = 256 // 默认写请求缓存容量
)

// Shard 表分片
type Shard struct {
	table      map[string]Value
	rw         sync.RWMutex
	writeQueue chan *writeReq // 写请求缓存队列
}

// 写请求
type writeReq struct {
	key   *Key
	value Value
}

func newShard() *Shard {
	table := &Shard{
		table:      make(map[string]Value),
		writeQueue: make(chan *writeReq, DefaultWriteQueueCap),
	}
	return table
}

func (t *Shard) Get(key *Key) Value {
	t.rw.RLock()
	defer t.rw.RUnlock()
	return t.table[key.GetKey()]
}

func (t *Shard) Set(key *Key, val Value) {
	t.writeQueue <- &writeReq{
		key:   key,
		value: val,
	}
	fmt.Println("写入", *key, val)
}

// todo GC
func (t *Shard) Delete(key *Key) {
	t.rw.Lock()
	defer t.rw.Unlock()
	t.table[key.GetKey()].DeleteValue()
}
