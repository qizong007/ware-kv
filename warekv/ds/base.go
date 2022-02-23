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

func (b *Base) Size() int {
	return 5 * 8 // 字段数 * 大小
}
