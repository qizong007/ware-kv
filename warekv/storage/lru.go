package storage

import (
	"container/list"
)

type entry struct {
	key   *Key
	value Value
}

type LRUCache struct {
	maxBytes  int64
	usedBytes int64
	ll        *list.List
	cache     map[string]*list.Element
}

func New(maxBytes int64) *LRUCache {
	return &LRUCache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}
}

func (c *LRUCache) Get(key *Key) Value {
	if ele, ok := c.cache[key.GetKey()]; ok {
		c.ll.MoveToFront(ele) // fresh
		kv := ele.Value.(*entry)
		return kv.value
	}
	return nil
}

func (c *LRUCache) Set(key *Key, value Value) {
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
	// 检测是否内存超了
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.removeOldest()
	}
}

func (c *LRUCache) removeOldest() {
	last := c.ll.Back() // 取最旧的元素
	if last != nil {
		c.ll.Remove(last)
		kv := last.Value.(*entry)
		delete(c.cache, kv.key.GetKey())
		// todo: c.usedBytes
		// c.usedBytes -= int64(len(kv.key) + kv.value.Len())
	}
}

func (c *LRUCache) SetInTime(key *Key, val Value) {
	c.Set(key, val)
}

func (c *LRUCache) Delete(key *Key) {
	// TODO
}

func (c *LRUCache) Len() int {
	return c.ll.Len()
}
