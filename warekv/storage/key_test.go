package storage

import (
	"fmt"
	"testing"
	"time"
)

func TestKeyHashcode(t *testing.T) {
	key := MakeKey("Aa")
	fmt.Println(key.Hashcode())
	key.SetKey("BB")
	fmt.Println(key.Hashcode())
	key.SetKey("hello")
	fmt.Println(key.Hashcode())
	key.SetKey("祺总007")
	tm := time.Now()
	key.Hashcode()
	fmt.Println(time.Since(tm))
	fmt.Println(key.Hashcode())
}
