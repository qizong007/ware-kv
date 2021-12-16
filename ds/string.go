package ds

import (
	"fmt"
	"time"
	"ware-kv/storage"
)

type String struct {
	Base
	Val string
}

func (s *String) GetValue() interface{} {
	return s.Val
}

func (s *String) SetValue(val interface{}) {
	s.Val = val.(string)
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
		Val: val,
	}
}

func Value2String(val storage.Value) *String {
	return val.(*String)
}

func (s *String) GetLen() int {
	return len(s.Val)
}

func (s *String) Append(str string) {
	s.Val = s.Val + str
}

// GetRange TODO 取范围（左闭右开）
// start为左， == -1 参数无效
//  end  == -1 参数无效
func (s *String) GetRange(start int, end int) (string, error) {
	if start < 0 && end < 0 {
		return "", fmt.Errorf("Start  and  End  Need  Greater  Than  0")
	}
	if start >= s.GetLen() || end >= s.GetLen() {
		//return "", fmt.Errorf(util.ParamErr)
	}
	if start >= end {
		//return "", fmt.Errorf(util.ParamErr)
	}
	if start >= 0 && end >= 0 {
		return s.Val[start:end], nil
	}
	if start >= 0 {
		return s.Val[start:], nil
	}
	return s.Val[:end], nil
}

// TODO
// 	 expire
// 	 set if not exist
