package main

import (
	"fmt"
	"github.com/qizong007/ware-kv/camera"
	"github.com/qizong007/ware-kv/tracker"
	"github.com/qizong007/ware-kv/warekv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type WareOption struct {
	Port       string                   `yaml:"Port"`
	WareEngine *warekv.WareEngineOption `yaml:"WareEngine"`
	Tracker    *tracker.TrackerOption   `yaml:"Tracker"`
	Camera     *camera.CameraOption     `yaml:"Camera"`
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
		WareEngine: warekv.DefaultOption(),
		Tracker:    tracker.DefaultOption(),
		Camera:     camera.DefaultOption(),
	}
}
