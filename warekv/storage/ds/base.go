package ds

import (
	"time"
	"ware-kv/warekv/util"
)

type Base struct {
	Type       util.DSType
	CreateTime int64
	UpdateTime int64
	DeleteTime int64
	ExpireTime *int64
	Version    uint64
}

func NewBase(tp util.DSType) *Base {
	t := time.Now().Unix()
	return &Base{
		Type:       tp,
		CreateTime: t,
		UpdateTime: t,
		DeleteTime: 0,
		ExpireTime: nil,
		Version:    1,
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

func (b *Base) Update() {
	b.UpdateTime = time.Now().Unix()
	b.Version++
}

func (b *Base) GetType() util.DSType {
	return b.Type
}

func (b *Base) GetBase() *Base {
	return b
}

func (b *Base) SetBase(base *Base) {
	b.Type = base.Type
	b.CreateTime = base.CreateTime
	b.UpdateTime = base.UpdateTime
	b.DeleteTime = base.DeleteTime
	b.ExpireTime = base.ExpireTime
	b.Version = base.Version
}
