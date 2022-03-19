package util

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestSkipList(t *testing.T) {
	sl := NewSkipList()
	sl.Insert(1, "forth")
	sl.Insert(-100, "second")
	sl.Insert(1, "third")
	sl.Insert(math.MinInt64, "first")
	sl.Insert(math.MaxInt64, "last")
	fmt.Println(sl.GetList())
	var ok bool
	nodes, ok := sl.Search(1)
	if ok {
		fmt.Println("found", nodes)
	} else {
		fmt.Println("not found", -100)
	}
	sl.Delete(100)
	fmt.Println(sl.GetList())
	sl.Delete(1)
	fmt.Println(sl.GetList())
	nodes, ok = sl.Search(1)
	if ok {
		fmt.Println("found", nodes)
	} else {
		fmt.Println("not found", -100)
	}
	fmt.Println(sl.GetList())
}

func TestSkipList_Insert(t *testing.T) {
	num := 100000
	scores := make([]float64, num)
	for i := 0; i < num; i++ {
		scores[i] = rand.Float64()
	}
	sl := NewSkipList()
	s := time.Now()
	for i := 0; i < num; i++ {
		sl.Insert(scores[i], i)
	}
	fmt.Println(time.Since(s))
}
