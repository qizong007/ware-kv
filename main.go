package main

import (
	"fmt"
	"ware-kv/initializer"
)

const (
	Port = "7777"
)

func main() {
	r := initializer.Boot()
	r.Run(fmt.Sprintf(":%s", Port))
}
