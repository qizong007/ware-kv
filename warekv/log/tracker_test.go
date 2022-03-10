package log

import (
	"testing"
	"time"
)

func TestTracker(t *testing.T) {
	tracker := NewTracker(nil)
	tracker.Write([]byte("s k1 v1\n"))
	time.Sleep(time.Second * 2)
	tracker.Write([]byte("s k2 [1,2,3]\n"))
	time.Sleep(time.Second * 2)
	tracker.Write([]byte("d k1\n"))
	time.Sleep(time.Second * 2)
}
