package ds

import (
	"fmt"
	"github.com/qizong007/ware-kv/warekv/util"
	"sync"
	"unsafe"
)

// ZList sorted-list
type ZList struct {
	Base
	skipList *util.SkipList
	rw       sync.RWMutex
}

var zListStructMemUsage int

func init() {
	zListStructMemUsage = int(unsafe.Sizeof(ZList{}))
}

type ZElement struct {
	Val   interface{}
	Score float64
}

func (zl *ZList) GetValue() interface{} {
	zl.rw.RLock()
	defer zl.rw.RUnlock()
	val := zl.skipList.GetList()
	return val
}

func (zl *ZList) Size() int {
	size := zListStructMemUsage
	if zl.ExpireTime != nil {
		size += 8
	}
	zl.rw.RLock()
	defer zl.rw.RUnlock()
	if rSize := util.GetRealSizeOf(zl.skipList); rSize > 0 {
		size += rSize
	}
	return size
}

func MakeZList(list []util.SlElement) *ZList {
	sl := util.NewSkipList()
	for i := range list {
		sl.Insert(list[i].Score, list[i].Val)
	}
	return &ZList{
		Base:     *NewBase(util.ZListDS),
		skipList: sl,
	}
}

// GetListBetween Left-Close and Right-Open
func (zl *ZList) GetListBetween(left int, right int) ([]util.SlElement, error) {
	zl.rw.RLock()
	defer zl.rw.RUnlock()
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
	zl.rw.RLock()
	defer zl.rw.RUnlock()
	zlLen := zl.GetLen()
	if left < 0 || left >= zlLen {
		return nil, fmt.Errorf("array out of bounds")
	}
	list := zl.skipList.GetList()
	return list[left:], nil
}

func (zl *ZList) GetListEndAt(right int) ([]util.SlElement, error) {
	zl.rw.RLock()
	defer zl.rw.RUnlock()
	zlLen := zl.GetLen()
	if right < 0 || right >= zlLen {
		return nil, fmt.Errorf("array out of bounds")
	}
	list := zl.skipList.GetList()
	return list[:right+1], nil
}

func (zl *ZList) GetElementAt(pos int) (*util.SlElement, error) {
	zl.rw.RLock()
	defer zl.rw.RUnlock()
	if pos < 0 || pos >= zl.GetLen() {
		return nil, fmt.Errorf("pos out of bounds")
	}
	list := zl.skipList.GetList()
	return &list[pos], nil
}

func (zl *ZList) GetListInScore(min float64, max float64) ([]util.SlElement, error) {
	zl.rw.RLock()
	defer zl.rw.RUnlock()
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
	zl.rw.Lock()
	defer zl.rw.Unlock()
	zl.Update()
	for i := range list {
		zl.skipList.Insert(list[i].Score, list[i].Val)
	}
}

func (zl *ZList) RemoveScore(score float64) {
	zl.rw.Lock()
	defer zl.rw.Unlock()
	zl.Update()
	zl.skipList.Delete(score)
}

func (zl *ZList) RemoveScores(scores []float64) {
	zl.rw.Lock()
	defer zl.rw.Unlock()
	zl.Update()
	for i := range scores {
		zl.skipList.Delete(scores[i])
	}
}

func (zl *ZList) RemoveInScore(min float64, max float64) error {
	list, err := zl.GetListInScore(min, max)
	if err != nil {
		return err
	}
	zl.rw.Lock()
	defer zl.rw.Unlock()
	zl.Update()
	for i := range list {
		zl.skipList.Delete(list[i].Score)
	}
	return nil
}

func (zl *ZList) GetLen() int {
	zl.rw.RLock()
	defer zl.rw.RUnlock()
	return zl.skipList.Len()
}
