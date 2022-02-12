package main

import (
	"fmt"
	"ware-kv/warekv"
)

const (
	Port = "7777"
)

func main() {
	r := warekv.Boot()
	r.Run(fmt.Sprintf(":%s", Port))
}
