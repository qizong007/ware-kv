package storage

import (
	"container/list"
	"github.com/qizong007/ware-kv/warekv/util"
	"sync"
)

const (
	LFUStrategy = "lfu"
)

type lfuEntry struct {
	key      string
	value    interface{}
	freqNode *list.Element
}

type LFUCache struct {
	maxBytes  int64
	usedBytes int64
	cache     map[string]*lfuEntry
	freqList  *list.List // from low to high freq
	rw        sync.RWMutex
}

type freqListEntry struct {
	entries map[*lfuEntry]struct{}
	freq    int
}

func NewLFUCache(maxBytes int64) *LFUCache {
	return &LFUCache{
		maxBytes: maxBytes,
		cache:    make(map[string]*lfuEntry),
		freqList: list.New(),
	}
}

func (c *LFUCache) Get(key *Key) Value {
	c.rw.RLock()
	defer c.rw.RUnlock()
	if e, ok := c.cache[key.GetKey()]; ok {
		c.increment(e)
		return e.value.(Value)
	}
	return nil
}

func (c *LFUCache) Set(k *Key, value Value) {
	c.rw.Lock()
	defer c.rw.Unlock()
	key := k.GetKey()
	if e, ok := c.cache[key]; ok {
		// update
		c.usedBytes += int64(value.Size() - e.value.(Value).Size())
		e.value = value
		c.increment(e)
	} else {
		// insert
		newEntry := &lfuEntry{
			key:   key,
			value: value,
		}
		c.cache[key] = newEntry
		c.usedBytes += int64(len(key) + value.Size())
		c.increment(newEntry)
	}
	// memory usage check
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.evict()
	}
}

func (c *LFUCache) evict() {
	head := c.freqList.Front()
	entries := head.Value.(*freqListEntry).entries
	for entry := range entries {
		c.remove(entry)
		break
	}
}

func (c *LFUCache) increment(e *lfuEntry) {
	var (
		nextFreq int
		next     *list.Element
	)

	current := e.freqNode

	if current == nil {
		// new entry
		nextFreq = 1
		next = c.freqList.Front()
	} else {
		// move up
		nextFreq = current.Value.(*freqListEntry).freq + 1
		next = current.Next()
	}

	if next == nil || next.Value.(*freqListEntry).freq != nextFreq {
		// create a new list entry
		newNode := &freqListEntry{
			entries: make(map[*lfuEntry]struct{}),
			freq:    nextFreq,
		}
		if current != nil {
			next = c.freqList.InsertAfter(newNode, current)
		} else {
			next = c.freqList.PushFront(newNode)
		}
	}

	e.freqNode = next
	next.Value.(*freqListEntry).entries[e] = struct{}{}

	if current != nil {
		c.removeEntryFromFreqNode(current, e)
	}
}

func (c *LFUCache) removeEntryFromFreqNode(freqNode *list.Element, entry *lfuEntry) {
	entries := freqNode.Value.(*freqListEntry).entries
	delete(entries, entry)
	if len(entries) == 0 {
		c.freqList.Remove(freqNode)
	}
}

func (c *LFUCache) SetInTime(k *Key, value Value) {
	c.Set(k, value)
}

func (c *LFUCache) Delete(key *Key) {
	c.rw.Lock()
	defer c.rw.Unlock()
	if e, ok := c.cache[key.GetKey()]; ok {
		c.remove(e)
	}
}

func (c *LFUCache) remove(entry *lfuEntry) {
	c.removeKVFromCache(entry)
	c.removeEntryFromFreqNode(entry.freqNode, entry)
}

func (c *LFUCache) removeKVFromCache(entry *lfuEntry) {
	delete(c.cache, entry.key)
	c.usedBytes -= int64(len(entry.key) + entry.value.(Value).Size())
}

func (c *LFUCache) Close() {
	// compatible with the interface
}

func (c *LFUCache) KeyNum() int {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return len(c.cache)
}

func (c *LFUCache) Type() string {
	return "LFUCache"
}

func (c *LFUCache) View() []byte {
	data := make([]byte, 0)
	// table flag
	data = append(data, uint8(TableFlag))
	// keys num
	data = append(data, util.IntToBytes(len(c.cache))...)
	// kv pairs
	for k, kv := range c.cache {
		data = append(data, kvPairView(k, kv.value.(Value))...)
	}
	return data
}

func (c *LFUCache) GetFlag() Flag {
	return TableFlag
}

func (c *LFUCache) MemUsage() int64 {
	c.rw.RLock()
	sum := int64(0)
	for _, e := range c.cache {
		val := e.value.(Value)
		sum += int64(val.Size())
	}
	c.rw.RUnlock()

	c.rw.Lock()
	// refresh mem usage
	c.usedBytes = sum
	// memory usage check
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.evict()
	}
	c.rw.Unlock()

	return c.usedBytes
}
