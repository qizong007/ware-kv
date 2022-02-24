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
	sl.Insert(1, "third")
	sl.Insert(-100, "second")
	sl.Insert(math.MinInt64, "first")
	sl.Insert(math.MaxInt64, "forth")
	fmt.Println(sl.GetList())
	var ok bool
	var node *slNode
	node, ok = sl.Search(-100)
	if ok {
		fmt.Println("found", node.score, node.val)
	} else {
		fmt.Println("not found", -100)
	}
	node = sl.Delete(100)
	fmt.Println(node)
	node = sl.Delete(-100)
	fmt.Println(node)
	node, ok = sl.Search(-100)
	if ok {
		fmt.Println("found", node.score, node.val)
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
