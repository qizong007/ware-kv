package util

type Set struct {
	set map[interface{}]void
}

type void struct{}

func NewSet(list []interface{}) *Set {
	s := make(map[interface{}]void)
	for _, v := range list {
		s[v] = void{}
	}
	return &Set{
		set: s,
	}
}

func (s *Set) Size() int {
	return len(s.set)
}

func (s *Set) Get() []interface{} {
	list := make([]interface{}, s.Size())
	i := 0
	for e := range s.set {
		list[i] = e
		i++
	}
	return list
}

func (s *Set) Contains(e interface{}) bool {
	_, ok := s.set[e]
	return ok
}

func (s *Set) Add(e interface{}) {
	if !s.Contains(e) {
		s.set[e] = void{}
	}
}

func (s *Set) Remove(e interface{}) {
	if s.Contains(e) {
		delete(s.set, e)
	}
}

// Intersect 交集
func (s *Set) Intersect(another *Set) *Set {
	res := make(map[interface{}]void)
	if s.Size() <= another.Size() {
		for e := range s.set {
			if another.Contains(e) {
				res[e] = void{}
			}
		}
	} else {
		for e := range another.set {
			if s.Contains(e) {
				res[e] = void{}
			}
		}
	}
	return &Set{set: res}
}

// Union 并集
func (s *Set) Union(another *Set) *Set {
	res := make(map[interface{}]void)
	for e := range s.set {
		res[e] = void{}
	}
	for e := range another.set {
		if !s.Contains(e) {
			res[e] = void{}
		}
	}
	return &Set{set: res}
}

// Diff 差集
func (s *Set) Diff(another *Set) *Set {
	res := make(map[interface{}]void)
	for e := range s.set {
		if !another.Contains(e) {
			res[e] = void{}
		}
	}
	for e := range another.set {
		if !s.Contains(e) {
			res[e] = void{}
		}
	}
	return &Set{set: res}
}
