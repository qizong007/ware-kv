package cache

import (
	"container/list"
	"ware-kv/warekv/storage"
)

type LRUCache struct {
	maxBytes  int64 // 最大内存
	usedBytes int64 // 已用内存
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

func (c *LRUCache) Get(key *storage.Key) (storage.Value, bool) {
	if ele, ok := c.cache[key.GetKey()]; ok {
		c.ll.MoveToFront(ele) // 设为最新
		kv := ele.Value.(*entry)
		return kv.value, ok
	}
	return nil, false
}

func (c *LRUCache) Set(key *storage.Key, value storage.Value) {
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

func (c *LRUCache) Delete(key *storage.Key) {
	// TODO
}

func (c *LRUCache) Len() int {
	return c.ll.Len()
}
