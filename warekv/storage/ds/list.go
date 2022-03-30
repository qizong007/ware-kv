package ds

import (
	"fmt"
	"github.com/qizong007/ware-kv/warekv/util"
	"reflect"
	"sync"
)

type List struct {
	Base
	list *[]interface{}
	rw   sync.RWMutex
}

func (l *List) GetValue() interface{} {
	val := l.listView()
	return val
}

// 深拷贝
func (l *List) listView() []interface{} {
	l.rw.RLock()
	defer l.rw.RUnlock()
	list := make([]interface{}, len(*l.list))
	for i, e := range *l.list {
		list[i] = e
	}
	return list
}

func MakeList(list []interface{}) *List {
	return &List{
		Base: *NewBase(util.ListDS),
		list: &list,
	}
}

func (l *List) GetListBetween(left int, right int) ([]interface{}, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()
	lLen := l.GetLen()
	if left < 0 || left >= lLen || right < 0 {
		return nil, fmt.Errorf("array out of bounds")
	}
	if left >= right {
		return nil, fmt.Errorf("left should not be larger than(or equal to) right")
	}
	if right >= lLen {
		right = lLen
	}
	list := l.listView()
	return list[left:right], nil
}

func (l *List) GetListStartWith(left int) ([]interface{}, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()
	lLen := l.GetLen()
	if left < 0 || left >= lLen {
		return nil, fmt.Errorf("array out of bounds")
	}
	list := l.listView()
	return list[left:], nil
}

func (l *List) GetListEndAt(right int) ([]interface{}, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()
	lLen := l.GetLen()
	if right < 0 || right >= lLen {
		return nil, fmt.Errorf("array out of bounds")
	}
	list := l.listView()
	return list[:right+1], nil
}

func (l *List) GetElementAt(pos int) (interface{}, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()
	if pos < 0 || pos >= l.GetLen() {
		return nil, fmt.Errorf("pos out of bounds")
	}
	list := l.listView()
	return list[pos], nil
}

func (l *List) GetListEqualToVal(val interface{}) []interface{} {
	list := make([]interface{}, 0)
	l.rw.RLock()
	defer l.rw.RUnlock()
	for _, v := range *l.list {
		if v == val {
			list = append(list, v)
		}
	}
	return list
}

func (l *List) Append(list []interface{}) {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.Update()
	*l.list = append(*l.list, list...)
}

func (l *List) RemoveAt(idx int) {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.Update()
	for i := range *l.list {
		if i == idx {
			*l.list = append((*l.list)[:i], (*l.list)[i+1:]...)
			break
		}
	}
}

func (l *List) RemoveVal(val interface{}) {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.Update()
	for i, v := range *l.list {
		if reflect.DeepEqual(v, val) {
			*l.list = append((*l.list)[:i], (*l.list)[i+1:]...)
		}
	}
}

func (l *List) GetLen() int {
	l.rw.RLock()
	defer l.rw.RUnlock()
	return len(*l.list)
}

func (l *List) RPush(ele interface{}) {
	l.Append([]interface{}{ele})
}

func (l *List) RPop() interface{} {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.Update()
	last := len(*l.list) - 1
	res := (*l.list)[last]
	*l.list = (*l.list)[:last]
	return res
}

func (l *List) LPush(ele interface{}) {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.Update()
	*l.list = append([]interface{}{ele}, *l.list...)
}

func (l *List) LPop() interface{} {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.Update()
	res := (*l.list)[0]
	*l.list = (*l.list)[1:]
	return res
}
