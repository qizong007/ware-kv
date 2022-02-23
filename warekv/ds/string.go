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

func (s *String) SetValue(val interface{}) {
	s.str = val.(string)
}

func MakeString(val string) *String {
	return &String{
		Base: *NewBase(),
		str: val,
	}
}

func Value2String(val storage.Value) *String {
	return val.(*String)
}

func (s *String) GetLen() int {
	return len(s.str)
}

func (s *String) Size() int {
	return s.Base.Size() + len(s.str)
}
