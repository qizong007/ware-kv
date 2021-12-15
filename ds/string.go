package ds

import (
	"fmt"
	"time"
)

type String struct {
	Val string
	CreateTime int64
	UpdateTime int64
	DeleteTime int64
}

func (s *String) GetValue() interface{} {
	return s.Val
}

func (s *String) SetValue(val interface{}) {
	s.Val = val.(string)
}

func (s *String) DeleteValue() {
	fmt.Println(s.DeleteTime)
	fmt.Println(s)
	s.DeleteTime = time.Now().Unix()
	fmt.Println(s.DeleteTime)
	fmt.Println(s)
}

func (s *String) IsAlive() bool {
	if s.DeleteTime == 0 {
		return true
	}
	return false
}

func MakeString(val string) *String {
	t := time.Now().Unix()
	return &String{
		Val:        val,
		CreateTime: t,
		UpdateTime: t,
		DeleteTime: 0,
	}
}

func (s *String) GetLen() int {
	return len(s.Val)
}

func (s *String) Append(str string) {
	s.Val = s.Val + str
}

// GetRange 取范围（左闭右开）
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
//   append
// 	 expire
// 	 set if not exist
