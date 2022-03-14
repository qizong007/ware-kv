package tracker

import (
	"testing"
	"time"
)

func TestTrackerWrite(t *testing.T) {
	option := &TrackerOption{
		FilePath:               "../tracker.log",
		BufRefreshTickInterval: defaultBufferTickInterval,
	}
	tk := NewTracker(option)
	c := NewCreateCommand("k1", StringStruct, "hello, ware-kv", time.Now().Unix(), 0)
	tk.Write(c)
	d := NewDeleteCommand("k1")
	tk.Write(d)
	cl := NewCreateCommand("k2", ListStruct, []interface{}{1, 3.14, "wq"}, time.Now().Unix(), 0)
	tk.Write(cl)
	time.Sleep(time.Second * 2)
}
