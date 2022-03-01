package ds

import (
	"fmt"
	"reflect"
	"sync"
	"ware-kv/warekv/storage"
)

type List struct {
	Base
	list *[]interface{}
	rw       sync.RWMutex
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

func (l *List) SetValue(val interface{}) {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.list = val.(*[]interface{})
}

func MakeList(list []interface{}) *List {
	return &List{
		Base: *NewBase(),
		list: &list,
	}
}

func Value2List(val storage.Value) *List {
	return val.(*List)
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
	*l.list = append(*l.list, list...)
}

func (l *List) RemoveAt(idx int) {
	l.rw.Lock()
	defer l.rw.Unlock()
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
