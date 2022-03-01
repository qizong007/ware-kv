package ds

import (
	"sync"
	"ware-kv/warekv/storage"
	"ware-kv/warekv/util"
)

type Set struct {
	Base
	set *util.Set
	rw  sync.RWMutex
}

func (s *Set) GetValue() interface{} {
	val := s.setView()
	return val
}

// 深拷贝
func (s *Set) setView() []interface{} {
	s.rw.RLock()
	defer s.rw.RUnlock()
	list := s.set.Get()
	return list
}

func MakeSet(list []interface{}) *Set {
	return &Set{
		Base: *NewBase(),
		set:  util.NewSet(list),
	}
}

func Value2Set(val storage.Value) *Set {
	return val.(*Set)
}

func (s *Set) GetSize() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.set.Size()
}

func (s *Set) Contains(e interface{}) bool {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.set.Contains(e)
}

func (s *Set) Add(e interface{}) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.set.Add(e)
}

func (s *Set) Remove(e interface{}) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.set.Remove(e)
}

// Intersect 交集
func (s *Set) Intersect(another *Set) *Set {
	s.rw.RLock()
	defer s.rw.RUnlock()
	another.rw.RLock()
	defer another.rw.RUnlock()
	list := s.set.Intersect(another.set).Get()
	return MakeSet(list)
}

// Union 并集
func (s *Set) Union(another *Set) *Set {
	s.rw.RLock()
	defer s.rw.RUnlock()
	another.rw.RLock()
	defer another.rw.RUnlock()
	list := s.set.Union(another.set).Get()
	return MakeSet(list)
}

// Diff 差集
func (s *Set) Diff(another *Set) *Set {
	s.rw.RLock()
	defer s.rw.RUnlock()
	another.rw.RLock()
	defer another.rw.RUnlock()
	list := s.set.Diff(another.set).Get()
	return MakeSet(list)
}
