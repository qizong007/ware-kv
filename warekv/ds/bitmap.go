package ds

import (
	"fmt"
	"sync"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
)

type Bitmap struct {
	Base
	bitmap *util.Bitmap
	rw     sync.RWMutex
}

func (b *Bitmap) GetValue() interface{} {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.bitmap.Value()
}

func MakeBitmap() *Bitmap {
	return &Bitmap{
		Base:   *NewBase(),
		bitmap: util.NewBitmap(),
	}
}

func Value2Bitmap(val storage.Value) *Bitmap {
	return val.(*Bitmap)
}

func (b *Bitmap) GetLen() int {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.bitmap.Len()
}

func (b *Bitmap) GetBitCount(start, end int) (int, error) {
	b.rw.RLock()
	defer b.rw.RUnlock()
	count := b.bitmap.BitCount(start, end)
	if count == -1 {
		return 0, fmt.Errorf("out of bitmap's boundary")
	}
	return count, nil
}

func (b *Bitmap) GetBit(num int) int {
	b.rw.RLock()
	defer b.rw.RUnlock()
	if b.bitmap.Has(num) {
		return 1
	}
	return 0
}

func (b *Bitmap) SetBit(num int) {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.bitmap.Set(num)
}

func (b *Bitmap) ClearBit(num int) {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.bitmap.Clear(num)
}