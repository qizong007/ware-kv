package cache

import (
	"container/list"
	"sync"
	"ware-kv/warekv/storage"
)

type entry struct {
	key   *storage.Key
	value storage.Value
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

func (c *LRUCache) Get(key *storage.Key) storage.Value {
	c.rw.RLock()
	defer c.rw.RUnlock()
	if ele, ok := c.cache[key.GetKey()]; ok {
		c.ll.MoveToFront(ele) // fresh
		kv := ele.Value.(*entry)
		return kv.value
	}
	return nil
}

func (c *LRUCache) Set(key *storage.Key, value storage.Value) {
	c.rw.Lock()
	defer c.rw.Unlock()
	if e, ok := c.cache[key.GetKey()]; ok {
		// update
		c.ll.MoveToFront(e)
		kv := e.Value.(*entry)
		// todo: c.usedBytes
		// c.usedBytes += int64(value.Len() - kv.value.Len())
		kv.value = value
	} else {
		// insert
		newEle := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key.GetKey()] = newEle
		// todo: c.usedBytes
		// c.usedBytes += int64(len(key) + value.Len())
	}
	// memory usage check
	//for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
	//	c.removeOldest()
	//}
}

func (c *LRUCache) removeOldest() {
	last := c.ll.Back() // get the oldest element
	if last != nil {
		c.ll.Remove(last)
		kv := last.Value.(*entry)
		delete(c.cache, kv.key.GetKey())
		// todo: c.usedBytes
		// c.usedBytes -= int64(len(kv.key) + kv.value.Len())
	}
}

func (c *LRUCache) SetInTime(key *storage.Key, val storage.Value) {
	c.Set(key, val)
}

func (c *LRUCache) Delete(key *storage.Key) {
	// TODO
}

func (c *LRUCache) Len() int {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.ll.Len()
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
