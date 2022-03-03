package util

import (
	"encoding/binary"
	"math"
)

type BloomFilter struct {
	n   uint64
	m   uint64
	k   uint64
	mem *Bitmap
}

func max(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}

func NewBloomFilter(m, k uint64) *BloomFilter {
	return &BloomFilter{
		m:   max(1, m),
		k:   max(1, k),
		mem: NewBitmapWithCap(m),
	}
}

func baseHashes(data []byte) [4]uint64 {
	var d digest128 // murmur hashing
	hash1, hash2, hash3, hash4 := d.sum256(data)
	return [4]uint64{
		hash1, hash2, hash3, hash4,
	}
}

func location(h [4]uint64, i uint) uint64 {
	ii := uint64(i)
	return h[ii%2] + ii*h[2+(((ii+(ii%2))%4)/2)]
}

func (f *BloomFilter) location(h [4]uint64, i uint) uint {
	return uint(location(h, i) % f.m)
}

func estimateParameters(n uint, p float64) (m uint64, k uint64) {
	m = uint64(math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
	k = uint64(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	return
}

func NewBloomFilterWithEstimates(n uint, fp float64) *BloomFilter {
	m, k := estimateParameters(n, fp)
	return NewBloomFilter(m, k)
}

func (f *BloomFilter) Add(data string) {
	h := baseHashes([]byte(data))
	for i := uint64(0); i < f.k; i++ {
		f.mem.Set(int(f.location(h, uint(i))))
	}
	f.n++
}

// Test 查看 data 是否存在
func (f *BloomFilter) Test(data string) bool {
	h := baseHashes([]byte(data))
	for i := uint64(0); i < f.k; i++ {
		if !f.mem.Has(int(f.location(h, uint(i)))) {
			return false
		}
	}
	return true
}

func (f *BloomFilter) ClearAll() *BloomFilter {
	f.mem = NewBitmapWithCap(f.m)
	f.n = 0
	return f
}

// EstimateFalsePositiveRate 估计 BloomFilter 的假阳性概率
func (f *BloomFilter) EstimateFalsePositiveRate(n uint) (fpRate float64) {
	rounds := uint32(100000)
	f.ClearAll()
	n1 := make([]byte, 4)
	for i := uint32(0); i < uint32(n); i++ {
		binary.BigEndian.PutUint32(n1, i)
		f.Add(string(n1))
	}
	fp := 0
	for i := uint32(0); i < rounds; i++ {
		binary.BigEndian.PutUint32(n1, i+uint32(n)+1)
		if f.Test(string(n1)) {
			fp++
		}
	}
	fpRate = float64(fp) / (float64(rounds))
	f.ClearAll()
	return
}

func (f *BloomFilter) Size() uint64 {
	return f.n
}

type BloomView struct {
	N   uint64
	M   uint64
	K   uint64
	Mem []uint64
}

func (f *BloomFilter) Value() *BloomView {
	return &BloomView{
		N:   f.n,
		M:   f.m,
		K:   f.k,
		Mem: f.mem.Value(),
	}
}
