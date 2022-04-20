package storage

import (
	"container/list"
	"github.com/qizong007/ware-kv/warekv/util"
	"sync"
)

const (
	LRUStrategy = "lru"
)

type lruEntry struct {
	key   *Key
	value Value
}

type LRUCache struct {
	maxBytes  int64
	usedBytes int64
	ll        *list.List
	cache     map[string]*list.Element
	rw        sync.RWMutex
}

func NewLRUCache(maxBytes int64) *LRUCache {
	return &LRUCache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}
}

func (c *LRUCache) Get(key *Key) Value {
	c.rw.RLock()
	defer c.rw.RUnlock()
	if ele, ok := c.cache[key.GetKey()]; ok {
		c.ll.MoveToFront(ele) // fresh
		kv := ele.Value.(*lruEntry)
		return kv.value
	}
	return nil
}

func (c *LRUCache) Set(key *Key, value Value) {
	c.rw.Lock()
	defer c.rw.Unlock()
	if e, ok := c.cache[key.GetKey()]; ok {
		// update
		c.ll.MoveToFront(e)
		kv := e.Value.(*lruEntry)
		c.usedBytes += int64(value.Size() - kv.value.Size())
		kv.value = value
	} else {
		// insert
		newEle := c.ll.PushFront(&lruEntry{
			key:   key,
			value: value,
		})
		c.cache[key.GetKey()] = newEle
		c.usedBytes += int64(len(key.GetKey()) + value.Size())
	}
	// memory usage check
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.removeOldest()
	}
}

func (c *LRUCache) remove(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*lruEntry)
	delete(c.cache, kv.key.GetKey())
	c.usedBytes -= int64(len(kv.key.GetKey()) + kv.value.Size())
}

func (c *LRUCache) removeOldest() {
	last := c.ll.Back() // get the oldest element
	if last != nil {
		c.remove(last)
	}
}

func (c *LRUCache) SetInTime(key *Key, val Value) {
	c.Set(key, val)
}

func (c *LRUCache) Delete(key *Key) {
	c.rw.Lock()
	defer c.rw.Unlock()
	if e, ok := c.cache[key.GetKey()]; ok {
		c.remove(e)
	}
}

func (c *LRUCache) DeleteInTime(key *Key) {
	c.Delete(key)
}

func (c *LRUCache) Close() {
	// compatible with the interface
}

func (c *LRUCache) KeyNum() int {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return len(c.cache)
}

func (c *LRUCache) Type() string {
	return "LRUCache"
}

func (c *LRUCache) View() []byte {
	data := make([]byte, 0)
	// table flag
	data = append(data, uint8(TableFlag))
	// keys num
	data = append(data, util.IntToBytes(len(c.cache))...)
	// kv pairs
	for k, e := range c.cache {
		kv := e.Value.(*lruEntry)
		data = append(data, kvPairView(k, kv.value)...)
	}
	return data
}

func (c *LRUCache) GetFlag() Flag {
	return TableFlag
}

func (c *LRUCache) MemUsage() int64 {
	c.rw.RLock()
	sum := int64(0)
	for _, e := range c.cache {
		kv := e.Value.(*lruEntry)
		sum += int64(kv.value.Size())
	}
	c.rw.RUnlock()

	c.rw.Lock()
	// refresh mem usage
	c.usedBytes = sum
	// memory usage check
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.removeOldest()
	}
	c.rw.Unlock()

	return c.usedBytes
}
