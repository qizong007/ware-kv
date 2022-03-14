package ds

import (
	"fmt"
	"testing"
)

func TestSetOp(t *testing.T) {
	set1 := MakeSet([]interface{}{1,3.14,"wq",false})
	set2 := MakeSet([]interface{}{2,3.14,"wq",true})
	fmt.Println(set1.Intersect(set2).GetValue())
	fmt.Println(set1.Union(set2).GetValue())
	fmt.Println(set1.Diff(set2).GetValue())
}
