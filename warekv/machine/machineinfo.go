package machine

import (
	"github.com/qizong007/ware-kv/warekv/storage"
	"github.com/qizong007/ware-kv/warekv/util"
	"github.com/shirou/gopsutil/process"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	defaultInfoFreshFrequency = 1000
	infoFreshFrequencyMin     = 100
	infoFreshFrequencyMax     = 5000
)

type Info struct {
	pid        int
	cpuPercent float64
	memPercent float32
	memAlloc   uint64
	memUsed    uint64
	infoTicker *time.Ticker
	closer     chan bool
}

var (
	wareInfo *Info
)

func NewWareInfo(option *WareInfoOption) *Info {
	infoFreshFrequency := time.Millisecond * time.Duration(defaultInfoFreshFrequency)
	if option != nil {
		freshFrequency := util.SetIfHitLimit(int(option.FreshFrequency), infoFreshFrequencyMin, infoFreshFrequencyMax)
		infoFreshFrequency = time.Millisecond * time.Duration(freshFrequency)
	}
	wareInfo = &Info{
		pid:        os.Getpid(),
		infoTicker: time.NewTicker(infoFreshFrequency),
		closer:     make(chan bool),
	}
	wareInfo.start()
	return wareInfo
}

type WareInfoOption struct {
	FreshFrequency uint `yaml:"FreshFrequency"`
}

func DefaultWareInfoOption() *WareInfoOption {
	return &WareInfoOption{FreshFrequency: defaultInfoFreshFrequency}
}

type InfoView struct {
	Pid        int     `json:"pid"`
	CpuPercent float64 `json:"cpu_percent"`
	MemPercent float32 `json:"mem_percent"`
	MemAlloc   uint64  `json:"mem_alloc"`
	MemUsed    uint64  `json:"mem_used"`
	KeysTotal  int     `json:"keys_total"`
	TableType  string  `json:"table_type"`
}

func GetWareInfo() *InfoView {
	return &InfoView{
		Pid:        wareInfo.pid,
		CpuPercent: wareInfo.cpuPercent,
		MemPercent: wareInfo.memPercent,
		MemAlloc:   wareInfo.memAlloc,
		MemUsed:    wareInfo.memUsed,
		KeysTotal:  storage.GlobalTable.KeyNum(),
		TableType:  storage.GlobalTable.Type(),
	}
}

func (i *Info) start() {
	wareInfo.updateInfo()
	go wareInfo.refresh()
	log.Println("MachineInfo's Refresh worker starts working...")
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
	i.memUsed = uint64(storage.GlobalTable.MemUsage())
}
