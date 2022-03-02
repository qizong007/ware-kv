package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBitmap(t *testing.T) {
	bitmap := NewBitmap()
	bitmap.Set(0)
	bitmap.Set(3)
	bitmap.Set(63)
	bitmap.Set(64)
	bitmap.Set(127)
	bitmap.Set(128)
	bitmap.Set(255)
	bitmap.Set(511)
	bitmap.Set(513)
	assert.Equal(t, "[0,3,63,64,127,128,255,511,513]", bitmap.String())
	bitmap.Clear(511)
	assert.Equal(t, false, bitmap.Has(511))
	assert.Equal(t, 8, bitmap.Len())
	assert.Equal(t, 7, bitmap.BitCount(3, 514))
}
