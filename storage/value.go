package storage

// Value 所有数据结构来实现它
type Value interface {
	SetValue(val interface{})
	GetValue() interface{}
}
