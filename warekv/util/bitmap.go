package util

import (
	"encoding/json"
)

// For this 'Bitmap', every single bit is born to be 0.

type Bitmap struct {
	val    []uint64
	length int
}

func NewBitmap() *Bitmap {
	return &Bitmap{}
}

// NewBitmapWithCap cap's unit is [bit]
func NewBitmapWithCap(cap uint64) *Bitmap {
	return &Bitmap{
		val: make([]uint64, 0, (cap>>6)+1),
	}
}

func getSegAndOffsetByNum(num int) (seg int, offset uint) {
	seg = num >> 6                  // num/64
	offset = uint(num - (seg << 6)) // num%64
	return
}

func (bitmap *Bitmap) getBit(seg int, offset uint) uint64 {
	return bitmap.val[seg] & (1 << offset)
}

func (bitmap *Bitmap) setBit(seg int, offset uint) {
	bitmap.val[seg] |= 1 << offset
	bitmap.length++
}

func (bitmap *Bitmap) clearBit(seg int, offset uint) {
	bitmap.val[seg] -= 1 << offset
	bitmap.length--
}

func (bitmap *Bitmap) Has(num int) bool {
	seg, offset := getSegAndOffsetByNum(num)
	return seg < len(bitmap.val) && (bitmap.getBit(seg, offset)) != 0
}

func (bitmap *Bitmap) Set(num int) {
	seg, offset := getSegAndOffsetByNum(num)
	// expand capacity
	for seg >= len(bitmap.val) {
		bitmap.val = append(bitmap.val, 0)
	}
	// check if 'num' has already been in the bitmap
	if bitmap.getBit(seg, offset) == 0 {
		bitmap.setBit(seg, offset)
	}
}

func (bitmap *Bitmap) Clear(num int) {
	seg, offset := getSegAndOffsetByNum(num)
	if seg >= len(bitmap.val) || bitmap.getBit(seg, offset) == 0 {
		return
	}
	bitmap.clearBit(seg, offset)
}

func (bitmap *Bitmap) Len() int {
	return bitmap.length
}

func (bitmap *Bitmap) Cap() int {
	return len(bitmap.val) << 6
}

func hammingWeight(n uint64) int {
	res := 0
	for n != 0 {
		n = n & (n - 1)
		res++
	}
	return res
}

func (bitmap *Bitmap) BitCount(start int, end int) int {
	if start < 0 || end < 0 || start >= bitmap.Cap() {
		return -1
	}
	if end > bitmap.Cap() {
		end = bitmap.Cap()
	}
	startSeg, startOffset := getSegAndOffsetByNum(start)
	endSeg, endOffset := getSegAndOffsetByNum(end)
	count := 0
	for i := startSeg; i <= endSeg; i++ {
		count += hammingWeight(bitmap.val[i])
	}
	// the segment 'start' at
	for i := uint(0); i < startOffset; i++ {
		if bitmap.getBit(startSeg, i) != 0 {
			count--
		}
	}
	// the segment 'end' at
	for i := endOffset; i < 64; i++ {
		if bitmap.getBit(endSeg, i) != 0 {
			count--
		}
	}
	return count
}

func (bitmap *Bitmap) Value() []uint64 {
	arr := make([]uint64, bitmap.length)
	index := 0
	for seg, v := range bitmap.val {
		if v == 0 {
			continue
		}
		// offset
		for j := uint(0); j < 64; j++ {
			if v&(1<<j) != 0 {
				arr[index] = uint64(64*uint(seg) + j)
				index++
			}
		}
	}
	return arr
}

func (bitmap *Bitmap) String() string {
	bytes, _ := json.Marshal(bitmap.Value())
	return string(bytes)
}
