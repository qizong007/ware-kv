package ds

import (
	"github.com/qizong007/ware-kv/warekv/util"
	"unsafe"
)

type String struct {
	Base
	str string
}

var stringStructMemUsage int

func init() {
	stringStructMemUsage = int(unsafe.Sizeof(String{}))
}

func (s *String) GetValue() interface{} {
	return s.str
}

func (s *String) Size() int {
	size := stringStructMemUsage
	if s.ExpireTime != nil {
		size += 8
	}
	size += len(s.str)
	return size
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
