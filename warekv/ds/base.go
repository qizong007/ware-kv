package ds

import (
	"time"
)

type Base struct {
	CreateTime int64
	UpdateTime int64
	DeleteTime int64
	ExpireTime *int64
	Version    int64
}

func NewBase() *Base {
	t := time.Now().Unix()
	return &Base{
		CreateTime: t,
		UpdateTime: t,
		DeleteTime: 0,
		ExpireTime: nil,
	}
}

func (b *Base) DeleteValue() {
	b.DeleteTime = time.Now().Unix()
}

func (b *Base) IsAlive() bool {
	if b.DeleteTime == 0 {
		return true
	}
	return false
}

func (b *Base) IsExpired() bool {
	if b.ExpireTime != nil && *b.ExpireTime <= time.Now().Unix() {
		return true
	}
	return false
}

func (b *Base) WithExpireTime(delta int64) {
	t := delta + time.Now().Unix()
	b.ExpireTime = &t
}
