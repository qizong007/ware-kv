package util

import (
	"fmt"
	"testing"
)

func TestBasicBloom(t *testing.T) {
	f := NewBloomFilter(1000, 4)
	n1 := "Bess"
	n2 := "Jane"
	n3 := "Emma"
	f.Add(n1)
	n1b := f.Test(n1)
	n2b := f.Test(n2)
	n3b := f.Test(n3)
	fmt.Println(n1b)
	fmt.Println(n2b)
	fmt.Println(n3b)
	fmt.Println(EstimateParameters(1000, 0.01))
}
