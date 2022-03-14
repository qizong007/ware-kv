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

// Deep-Copy
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
	s.Update()
	s.set.Add(e)
}

func (s *Set) Remove(e interface{}) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.Update()
	s.set.Remove(e)
}

func getTwoSetView(s1 *Set, s2 *Set) (*util.Set, *util.Set) {
	s1.rw.RLock()
	list1 := s1.set.Get()
	s1.rw.RUnlock()
	s2.rw.RLock()
	list2 := s2.set.Get()
	s2.rw.RUnlock()
	return util.NewSet(list1), util.NewSet(list2)
}

func (s *Set) Intersect(another *Set) *Set {
	set1, set2 := getTwoSetView(s, another)
	list := set1.Intersect(set2).Get()
	return MakeSet(list)
}

func (s *Set) Union(another *Set) *Set {
	set1, set2 := getTwoSetView(s, another)
	list := set1.Union(set2).Get()
	return MakeSet(list)
}

func (s *Set) Diff(another *Set) *Set {
	set1, set2 := getTwoSetView(s, another)
	list := set1.Diff(set2).Get()
	return MakeSet(list)
}
