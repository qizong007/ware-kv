package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/qizong007/ware-kv/util"
	"log"
	"os"
)

const (
	defaultLoadOptionPath = "ware.yaml"
	versionInfo           = "Ware-kv is now version 0.0.1"
)

var (
	help       bool
	v          bool
	t          bool
	configPath string
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.BoolVar(&t, "t", false, "test configuration and exit")
	flag.StringVar(&configPath, "c", defaultLoadOptionPath, "set configuration `file`")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `ware-kv version: ware-kv/%s
Usage: ware-kv [-hvt] [-c filename]

Options:
`, util.WareKVVersion)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	if v {
		fmt.Println(versionInfo)
		return
	}

	optionPath := defaultLoadOptionPath
	if configPath != "" {
		optionPath = configPath
	}

	var option *WareOption
	if _, err := os.Stat(optionPath); err == nil {
		// file exist
		option, err = LoadOption(optionPath)
		if err != nil {
			log.Fatalln("Load option file failed! Err:", err)
			return
		}
		if t {
			return
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// file not exist
		fmt.Printf("%q is not exists...\n", optionPath)
		if t {
			return
		}
		fmt.Println("This will load the default option...")
		option = DefaultOption()
	} else {
		fmt.Println("Load option file failed! Err:", err)
		return
	}

	if option == nil {
		panic("No Option!!! Please check your 'ware.yml'!!!")
	}
	Boot(option)
}
