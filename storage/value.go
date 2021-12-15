package storage

// Value 所有数据结构都需要实现它
type Value interface {
	SetValue(val interface{})
	GetValue() interface{}
	DeleteValue()
	IsAlive() bool
}
