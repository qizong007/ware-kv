package ds

import (
	"ware-kv/warekv/util"
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
		Base: *NewBase(util.StringDS),
		str:  val,
	}
}

func (s *String) GetLen() int {
	return len(s.str)
}
