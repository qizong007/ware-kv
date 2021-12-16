package ds

import "time"

type Base struct {
	CreateTime int64
	UpdateTime int64
	DeleteTime int64
	ExpireTime *int64
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
