package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"ware-kv/warekv"
)

type WareOption struct {
	Port       string                   `yaml:"Port"`
	WareEngine *warekv.WareEngineOption `yaml:"WareEngine"`
}

func LoadOption(file string) (*WareOption, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var conf WareOption
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", file, err)
	}
	return &conf, nil
}

func DefaultOption() *WareOption {
	return &WareOption{
		Port:       defaultPort,
		WareEngine: nil,
	}
}
