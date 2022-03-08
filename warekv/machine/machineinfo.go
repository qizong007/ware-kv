package machine

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	defaultInfoFreshFrequency = 1000
)

type Info struct {
	pid        int
	cpuPercent float64
	memPercent float32
	memAlloc   uint64
	infoTicker *time.Ticker
	closer     chan bool
}

var (
	wareInfo *Info
)

func NewWareInfo(option *WareInfoOption) *Info {
	infoFreshFrequency := time.Millisecond * time.Duration(defaultInfoFreshFrequency)
	if option != nil {
		infoFreshFrequency = time.Millisecond * time.Duration(option.FreshFrequency)
	}
	wareInfo = &Info{
		pid:        os.Getpid(),
		infoTicker: time.NewTicker(infoFreshFrequency),
		closer:     make(chan bool),
	}
	return wareInfo
}

type WareInfoOption struct {
	FreshFrequency uint `yaml:"FreshFrequency"`
}

func DefaultWareInfoOption() *WareInfoOption {
	return &WareInfoOption{FreshFrequency: defaultInfoFreshFrequency}
}

type InfoView struct {
	Pid        int
	CpuPercent float64
	MemPercent float32
	MemAlloc   uint64
}

func GetWareInfo() *InfoView {
	return &InfoView{
		Pid:        wareInfo.pid,
		CpuPercent: wareInfo.cpuPercent,
		MemPercent: wareInfo.memPercent,
		MemAlloc:   wareInfo.memAlloc,
	}
}

func (i *Info) Start() {
	wareInfo.updateInfo()
	go wareInfo.refresh()
	fmt.Println("MachineInfo's Refresh worker starts working...")
}

func (i *Info) Close() {
	i.closer <- true
}

func (i *Info) refresh() {
	for {
		select {
		case <-i.infoTicker.C:
			i.updateInfo()
		case <-i.closer:
			log.Println("MachineInfo's Refresh worker stops working...")
			return
		}
	}
}

func (i *Info) updateInfo() {
	processes, err := process.Processes()
	if err != nil {
		return
	}
	pid := int32(i.pid)
	for _, p := range processes {
		if p.Pid == pid {
			i.cpuPercent, _ = p.CPUPercent()
			i.memPercent, _ = p.MemoryPercent()
			break
		}
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	i.memAlloc = ms.Alloc
}
