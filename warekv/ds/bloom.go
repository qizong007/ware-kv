package ds

import (
	"sync"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
)

type BloomFilter struct {
	Base
	filter *util.BloomFilter
	rw     sync.RWMutex
}

func (b *BloomFilter) GetValue() interface{} {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.filter.Value()
}

// MakeBloomFilterSpecific 精确创建
func MakeBloomFilterSpecific(m, k uint64) *BloomFilter {
	return &BloomFilter{
		Base:   *NewBase(),
		filter: util.NewBloomFilter(m, k),
	}
}

// MakeBloomFilterFuzzy 模糊创建
func MakeBloomFilterFuzzy(n uint, fp float64) *BloomFilter {
	return &BloomFilter{
		Base:   *NewBase(),
		filter: util.NewBloomFilterWithEstimates(n, fp),
	}
}

func Value2BloomFilter(val storage.Value) *BloomFilter {
	return val.(*BloomFilter)
}

func (b *BloomFilter) GetSize() uint64 {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.filter.Size()
}

func (b *BloomFilter) Add(data string) {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.filter.Add(data)
}

func (b *BloomFilter) Test(data string) bool {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.filter.Test(data)
}

func (b *BloomFilter) ClearAll() {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.filter.ClearAll()
}

func (b *BloomFilter) EstimateFalsePositiveRate(n uint) float64 {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.filter.EstimateFalsePositiveRate(n)
}
