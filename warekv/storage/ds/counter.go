package ds

import (
	"github.com/qizong007/ware-kv/warekv/util"
	"sync/atomic"
)

type Counter struct {
	Base
	num int64
}

func (c *Counter) GetValue() interface{} {
	return atomic.LoadInt64(&c.num)
}

func MakeCounter(val int64) *Counter {
	return &Counter{
		Base: *NewBase(util.CounterDS),
		num:  val,
	}
}

func (c *Counter) IncrBy(delta int64) {
	c.Update()
	atomic.AddInt64(&c.num, delta)
}

func (c *Counter) DecrBy(delta int64) {
	c.IncrBy(-delta)
}

func (c *Counter) Incr() {
	c.IncrBy(1)
}

func (c *Counter) Decr() {
	c.IncrBy(-1)
}
