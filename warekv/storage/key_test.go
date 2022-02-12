package storage

import (
	"fmt"
	"testing"
	"time"
)

func TestKeyHashcode(t *testing.T) {
	a := "Aa"
	key := &Key{}
	key.Val = a
	fmt.Println(key.Hashcode())
	b := "BB"
	key = &Key{}
	key.Val = b
	fmt.Println(key.Hashcode())
	c := "hello"
	key = &Key{}
	key.Val = c
	fmt.Println(key.Hashcode())
	d := "祺总007"
	key = &Key{}
	key.Val = d
	tm := time.Now()
	key.Hashcode()
	fmt.Println(time.Since(tm))
	fmt.Println(key.Hashcode())
}
