package log

import (
	"testing"
	"time"
)

func TestTracker(t *testing.T) {
	tracker := NewTracker(nil)
	c := GenCreateCommand("k1", []interface{}{1,2.3,"asd",[]int{1,2,3}}, time.Now().Unix(), 1)
	tracker.Write(c)
	time.Sleep(time.Second * 2)
	d := GenDeleteCommand("k1")
	tracker.Write(d)
	time.Sleep(time.Second * 2)
}
