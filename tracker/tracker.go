package tracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	defaultTrackPath          = "./tracker.log"
	defaultBufferTickInterval = 1000
)

var tracker *Tracker

type Tracker struct {
	file       *os.File
	buffer     []byte
	bufLock    sync.Mutex
	bufTicker  *time.Ticker
	closer     chan bool
	isRealTime bool
}

type TrackerOption struct {
	FilePath               string `yaml:"FilePath"`
	BufRefreshTickInterval uint   `yaml:"BufRefreshTickInterval"`
}

func DefaultOption() *TrackerOption {
	return &TrackerOption{
		FilePath:               defaultTrackPath,
		BufRefreshTickInterval: defaultBufferTickInterval,
	}
}

func NewTracker(option *TrackerOption) *Tracker {
	filePath := defaultTrackPath
	bufTickInterval := uint(defaultBufferTickInterval)
	isRealTime := false
	if option != nil {
		filePath = option.FilePath
		bufTickInterval = option.BufRefreshTickInterval
		if bufTickInterval == 0 {
			isRealTime = true
		}
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("NewTracker Fail: %v", err))
		return nil
	}
	tracker = &Tracker{
		file:       file,
		buffer:     make([]byte, 0),
		isRealTime: isRealTime,
	}
	if !isRealTime {
		tracker.closer = make(chan bool)
		tracker.bufTicker = time.NewTicker(time.Duration(bufTickInterval) * time.Millisecond)
		tracker.start()
	} else {
		log.Println("Tracker start real-time refresh mode...")
	}
	return tracker
}

func GetTracker() *Tracker {
	return tracker
}

func (t *Tracker) start() {
	go t.scheduledRefresh()
	log.Println("Tracker's Refresh worker starts working...")
}

func (t *Tracker) Close() {
	if t.isRealTime {
		return
	}
	t.closer <- true
}

func (t *Tracker) LoadTracker() {
	log.Println("Tracker start loading...")
	start := time.Now()
	data, err := ioutil.ReadAll(t.file)
	if err != nil {
		panic(fmt.Sprintf("loadTracker Fail: %v", err))
		return
	}
	if len(data) == 0 {
		fmt.Println("Nothing in", t.file.Name())
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var command Command
		switch line[0:1] {
		case CreateOp:
			command = &CreateCommand{}
		case ModifyOp:
			command = &ModifyCommand{}
		case DeleteOp:
			command = &DeleteCommand{}
		case SubscribeOp:
			command = &SubCommand{}
		}
		resolveCommand(line[1:], command)
		command.Execute()
	}
	log.Printf("Tracker finish loading in %s...\n", time.Since(start).String())
}

func resolveCommand(command string, cmd Command) {
	err := json.Unmarshal([]byte(command), cmd)
	if err != nil {
		log.Println("genCreateCommand json.Unmarshal fail", err)
		return
	}
}

func (t *Tracker) scheduledRefresh() {
	for {
		select {
		case <-t.bufTicker.C:
			t.refresh()
		case <-t.closer:
			t.bufTicker.Stop()
			close(t.closer)
			t.buffer = nil
			_ = t.file.Close()
			log.Println("Tracker's Refresh worker starts working...")
			return
		}
	}
}

func (t *Tracker) Write(command Command) {
	t.bufLock.Lock()
	defer t.bufLock.Unlock()
	t.buffer = append(t.buffer, []byte(command.GetOpType()+command.String()+"\n")...)
	if t.isRealTime {
		t.flushToDisk()
	}
}

func (t *Tracker) refresh() {
	t.bufLock.Lock()
	defer t.bufLock.Unlock()
	t.flushToDisk()
}

func (t *Tracker) flushToDisk() {
	if len(t.buffer) == 0 {
		return
	}
	_, err := t.file.Write(t.buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.buffer = t.buffer[:0]
}
