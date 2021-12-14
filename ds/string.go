package ds

import (
	"fmt"
)

type String struct {
	Val string
}

func (s String) GetValue() interface{} {
	return s.Val
}

func (s String) SetValue(val interface{}) {
	s.Val = val.(string)
}

func (s String) GetLen() int {
	return len(s.Val)
}

func (s String) Append(str string) {
	s.Val = s.Val + str
}

// GetRange 取范围（左闭右开）
// start为左， == -1 参数无效
//  end  == -1 参数无效
func (s String) GetRange(start int, end int) (string, error) {
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
//   append
// 	 expire
// 	 set if not exist
