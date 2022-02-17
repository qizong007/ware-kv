package util

import (
	"fmt"
	"math"
	"testing"
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
