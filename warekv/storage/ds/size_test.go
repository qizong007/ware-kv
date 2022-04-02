package ds

import (
	"fmt"
	"github.com/qizong007/ware-kv/warekv/util"
	"github.com/rs/xid"
	"testing"
)

func TestBitMapSize(t *testing.T) {
	bm := MakeBitmap()
	bm.SetBit(1)
	bm.SetBit(2)
	bm.SetBit(7)
	bm.SetBit(8)
	fmt.Println(bm.Size())
}

func TestBloomFilterSize(t *testing.T) {
	bf := MakeBloomFilterFuzzy(BloomFilterFuzzyOption{
		N:  1000,
		Fp: 0.01,
	})
	fmt.Println(bf.Size())
}

func TestCounterSize(t *testing.T) {
	cnt := MakeCounter(10)
	fmt.Println(cnt.Size())
}

func TestListSize(t *testing.T) {
	list := MakeList([]interface{}{1, 3.14, "string", false, []int{1, 2, 3}, &Base{}})
	fmt.Println(list.Size())
}

func TestLockSize(t *testing.T) {
	lk := MakeLock()
	guid := xid.New().String()
	_ = lk.Lock(10, guid)
	fmt.Println(lk.Size())
}

func TestObjectSize(t *testing.T) {
	dict := map[string]interface{}{
		"name":     "wq",
		"age":      18,
		"handsome": true,
	}
	obj := MakeObject(dict)
	fmt.Println(obj.Size())
}

func TestSetSize(t *testing.T) {
	list := MakeSet([]interface{}{1, 3.1415, "string", false, &Base{}})
	fmt.Println(list.Size())
}

func TestStringSize(t *testing.T) {
	str := MakeString("warekv")
	fmt.Println(str.Size())
}

func TestZListSize(t *testing.T) {
	param := make([]util.SlElement, 0)
	param = append(param, util.SlElement{Score: 2, Val: 1})
	param = append(param, util.SlElement{Score: 1, Val: 3.14})
	param = append(param, util.SlElement{Score: 4, Val: "string"})
	param = append(param, util.SlElement{Score: 5, Val: false})
	param = append(param, util.SlElement{Score: 6, Val: []int{1, 2, 3}})
	param = append(param, util.SlElement{Score: 3, Val: &Base{}})
	list := MakeZList(param)
	fmt.Println(list.Size())
}
