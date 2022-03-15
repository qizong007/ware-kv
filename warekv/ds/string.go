package ds

import (
	"ware-kv/warekv/storage"
)

type String struct {
	Base
	str string
}

func (s *String) GetValue() interface{} {
	return s.str
}

func MakeString(val string) *String {
	return &String{
		Base: *NewBase(StringDS),
		str:  val,
	}
}

func Value2String(val storage.Value) *String {
	return val.(*String)
}

func (s *String) GetLen() int {
	return len(s.str)
}
