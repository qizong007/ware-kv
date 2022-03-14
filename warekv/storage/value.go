package storage

// Value all the data structure should implements it
type Value interface {
	GetValue() interface{}
	DeleteValue()
	IsAlive() bool
	IsExpired() bool
	WithExpireTime(t int64)
	Update()
}
