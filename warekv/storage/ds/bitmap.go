package ds

import (
	"fmt"
	"github.com/qizong007/ware-kv/warekv/util"
	"sync"
	"unsafe"
)

type Bitmap struct {
	Base
	bitmap *util.Bitmap
	rw     sync.RWMutex
}

var bitmapStructMemUsage int

func init() {
	bitmapStructMemUsage = int(unsafe.Sizeof(Bitmap{}))
}

func (b *Bitmap) GetValue() interface{} {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.bitmap.Value()
}

func (b *Bitmap) Size() int {
	baseSize := bitmapStructMemUsage
	if b.ExpireTime != nil {
		baseSize += 8
	}
	b.rw.RLock()
	defer b.rw.RUnlock()
	return baseSize + b.bitmap.MemoryUsage()
}

func MakeBitmap() *Bitmap {
	return &Bitmap{
		Base:   *NewBase(util.BitmapDS),
		bitmap: util.NewBitmap(),
	}
}

func MakeBitmapFromList(list []uint64) *Bitmap {
	bm := &Bitmap{
		Base:   *NewBase(util.BitmapDS),
		bitmap: util.NewBitmap(),
	}
	for i := range list {
		bm.bitmap.Set(int(list[i]))
	}
	return bm
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
	b.Update()
	b.bitmap.Set(num)
}

func (b *Bitmap) ClearBit(num int) {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.Update()
	b.bitmap.Clear(num)
}
