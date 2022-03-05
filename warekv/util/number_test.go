package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNearest2Power(t *testing.T) {
	assert.Equal(t, uint(8), Nearest2Power(7))
	assert.Equal(t, uint(8), Nearest2Power(8))
	assert.Equal(t, uint(0), Nearest2Power(0))
	assert.Equal(t, uint(1), Nearest2Power(1))
	assert.Equal(t, uint(2), Nearest2Power(2))
	assert.Equal(t, uint(4), Nearest2Power(3))
}
