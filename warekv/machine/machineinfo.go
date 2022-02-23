package machine

import (
	"github.com/shirou/gopsutil/process"
	"os"
	"runtime"
	"time"
)

const (
	defaultInfoFreshFrequency = time.Second
)

type Info struct {
	Pid        int
	CpuPercent float64
	MemPercent float32
	MemAlloc   uint64
}

var (
	wareInfo   *Info
	infoTicker *time.Ticker
)

func init() {
	wareInfo = &Info{
		Pid: os.Getpid(),
	}
	infoTicker = time.NewTicker(defaultInfoFreshFrequency)
	wareInfo.updateInfo()
	go wareInfo.refresh()
}

func GetWareInfo() *Info {
	return wareInfo
}

func (i *Info) refresh() {
	for {
		select {
		case <-infoTicker.C:
			i.updateInfo()
		}
	}
}

func (i *Info) updateInfo() {
	processes, err := process.Processes()
	if err != nil {
		return
	}
	pid := int32(i.Pid)
	for _, p := range processes {
		if p.Pid == pid {
			i.CpuPercent, _ = p.CPUPercent()
			i.MemPercent, _ = p.MemoryPercent()
			break
		}
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	i.MemAlloc = ms.Alloc
}
