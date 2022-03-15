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

type BloomFilterSpecificOption struct {
	M uint64
	K uint64
}

// MakeBloomFilterSpecific 精确创建
func MakeBloomFilterSpecific(option BloomFilterSpecificOption) *BloomFilter {
	return &BloomFilter{
		Base:   *NewBase(BloomFilterDS),
		filter: util.NewBloomFilter(option.M, option.K),
	}
}

type BloomFilterFuzzyOption struct {
	N  uint
	Fp float64
}

// MakeBloomFilterFuzzy 模糊创建
func MakeBloomFilterFuzzy(option BloomFilterFuzzyOption) *BloomFilter {
	return &BloomFilter{
		Base:   *NewBase(BloomFilterDS),
		filter: util.NewBloomFilterWithEstimates(option.N, option.Fp),
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
	b.Update()
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
	b.Update()
	b.filter.ClearAll()
}

func (b *BloomFilter) EstimateFalsePositiveRate(n uint) float64 {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.filter.EstimateFalsePositiveRate(n)
}
