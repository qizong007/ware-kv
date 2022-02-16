package ds

import (
	"time"
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
	t := time.Now().Unix()
	return &String{
		Base: Base{
			CreateTime: t,
			UpdateTime: t,
			DeleteTime: 0,
			ExpireTime: nil,
		},
		str: val,
	}
}

func Value2String(val storage.Value) *String {
	return val.(*String)
}

func (s *String) GetLen() int {
	return len(s.str)
}
