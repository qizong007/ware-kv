package storage

import (
	"sync"
	"time"
)

type writeEvent int

const (
	SetEvent = iota
	DeleteEvent
)

type Shard struct {
	table      map[string]Value
	rw         sync.RWMutex
	writeQueue chan *writeReq // write request's queue
	wqTicker   *time.Ticker   // write request's ticker
	gc         *WareGC
	closer     chan bool
}

type writeReq struct {
	event writeEvent
	key   *Key
	value Value
}

func newShard(writeQueueCap int, writeTickInterval time.Duration, gcOption *WareGCOption) *Shard {
	shard := &Shard{
		table:      make(map[string]Value),
		writeQueue: make(chan *writeReq, writeQueueCap),
		wqTicker:   time.NewTicker(writeTickInterval),
		closer:     make(chan bool),
	}
	shard.gc = NewWareGC(shard, gcOption)
	return shard
}

func (s *Shard) Start() {
	go s.scheduledBatchCommit()
}

func (s *Shard) Close() {
	s.closer <- true
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
		case <-s.wqTicker.C: // Batch WRITE
			if len(s.writeQueue) == 0 {
				continue
			}
			s.handleWriteQueue()
		case <-s.gc.gcTicker.C: // Batch SWEEP
			if len(s.gc.gcTasks) == 0 {
				continue
			}
			s.handleGC()
		case <-s.closer:
			s.gc.Close()
			s.closeWriteWorker()
			close(s.closer)
			return
		}
	}
}

func (s *Shard) handleWriteQueue() {
	s.rw.Lock()
	for entry := range s.writeQueue {
		key := entry.key.GetKey()
		switch entry.event {
		case SetEvent:
			s.table[key] = entry.value
		case DeleteEvent:
			s.table[key].DeleteValue()
			s.gc.Commit(key)
		}
		if len(s.writeQueue) == 0 {
			break
		}
	}
	s.rw.Unlock()
}

func (s *Shard) handleGC() {
	s.rw.Lock()
	for key := range s.gc.gcTasks {
		delete(s.table, key)
		if len(s.gc.gcTasks) == 0 {
			break
		}
	}
	s.rw.Unlock()
}

func (s *Shard) closeWriteWorker() {
	s.handleWriteQueue()
	close(s.writeQueue)
	s.wqTicker.Stop()
}

func (s *Shard) Count() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.table)
}
