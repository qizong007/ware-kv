package storage

import (
	"container/list"
	"fmt"
	"github.com/qizong007/ware-kv/warekv/storage/ds"
	"testing"
)

func printFreqList(ll *list.List) {
	n := ll.Len()
	cur := ll.Front()
	for i := 0; i < n; i++ {
		freqNode := cur.Value.(*freqListEntry)
		fmt.Print(*freqNode, " ")
		cur = cur.Next()
	}
	fmt.Println()
}

func TestLFU(t *testing.T) {
	cache := NewLFUCache(204)
	k1 := MakeKey("k1")
	k2 := MakeKey("k2")
	k3 := MakeKey("k3")
	k4 := MakeKey("k4")
	cache.Set(k1, ds.MakeString("v1"))
	cache.Set(k2, ds.MakeString("v2"))
	cache.Set(k3, ds.MakeString("v3"))
	cache.Set(k4, ds.MakeString("v4"))
	fmt.Println(cache.usedBytes)
	printFreqList(cache.freqList)
	fmt.Println(cache.Get(k1))
	fmt.Println(cache.Get(k2))
	fmt.Println(cache.Get(k2))
	fmt.Println(cache.Get(k3))
	fmt.Println(cache.Get(k3))
	fmt.Println(cache.Get(k3))
	fmt.Println(cache.Get(k4))
	cache.Delete(k4)
	fmt.Println(cache.usedBytes)
	printFreqList(cache.freqList)
}
