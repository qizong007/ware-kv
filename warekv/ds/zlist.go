package ds

import (
	"fmt"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
)

// ZList 有序列表
type ZList struct {
	Base
	skipList *util.SkipList
}

type ZElement struct {
	Val   interface{}
	Score float64
}

func (zl *ZList) GetValue() interface{} {
	return zl.skipList.GetList()
}

func (zl *ZList) SetValue(val interface{}) {
	zl.skipList = val.(*util.SkipList)
}

func MakeZList(list []util.SlElement) *ZList {
	sl := util.NewSkipList()
	for i := range list {
		sl.Insert(list[i].Score, list[i].Val)
	}
	return &ZList{
		Base:     *NewBase(),
		skipList: sl,
	}
}

func Value2ZList(val storage.Value) *ZList {
	return val.(*ZList)
}

// GetListBetween 左闭右开
func (zl *ZList) GetListBetween(left int, right int) ([]util.SlElement, error) {
	zlLen := zl.GetLen()
	if left < 0 || left >= zlLen || right < 0 {
		return nil, fmt.Errorf("array out of bounds")
	}
	if left >= right {
		return nil, fmt.Errorf("left should not be larger than(or equal to) right")
	}
	if right >= zlLen {
		right = zlLen
	}
	list := zl.skipList.GetList()
	return list[left:right], nil
}

func (zl *ZList) GetListStartWith(left int) ([]util.SlElement, error) {
	zlLen := zl.GetLen()
	if left < 0 || left >= zlLen {
		return nil, fmt.Errorf("array out of bounds")
	}
	list := zl.skipList.GetList()
	return list[left:], nil
}

func (zl *ZList) GetListEndAt(right int) ([]util.SlElement, error) {
	zlLen := zl.GetLen()
	if right < 0 || right >= zlLen {
		return nil, fmt.Errorf("array out of bounds")
	}
	list := zl.skipList.GetList()
	return list[:right], nil
}

func (zl *ZList) GetListInScore(min float64, max float64) ([]util.SlElement, error) {
	if min > max {
		return nil, fmt.Errorf("min should not be larger than max")
	}
	list := zl.skipList.GetList()
	res := make([]util.SlElement, 0)
	for i := range list {
		if list[i].Score > max {
			break
		}
		if list[i].Score >= min && list[i].Score <= max {
			res = append(res, list[i])
		}
	}
	return res, nil
}

func (zl *ZList) Add(list []util.SlElement) {
	for i := range list {
		zl.skipList.Insert(list[i].Score, list[i].Val)
	}
}

func (zl *ZList) Remove(score float64) {
	zl.skipList.Delete(score)
}

func (zl *ZList) RemoveInScore(min float64, max float64) error {
	list, err := zl.GetListInScore(min, max)
	if err != nil {
		return err
	}
	for i := range list {
		zl.skipList.Delete(list[i].Score)
	}
	return nil
}

func (zl *ZList) GetLen() int {
	return zl.skipList.Len()
}
