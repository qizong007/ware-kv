package log

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	defaultTrackPath          = "../../tracker.log"
	defaultBufferTickInterval = 1000
)

type Tracker struct {
	file      *os.File
	buffer    []byte
	bufLock   sync.Mutex
	bufTicker *time.Ticker
	closer    chan bool
}

type TrackerOption struct {
	FilePath               string
	BufRefreshTickInterval uint
}

type Command struct {
}

func NewTracker(option *TrackerOption) *Tracker {
	filePath := defaultTrackPath
	bufTickInterval := uint(defaultBufferTickInterval)
	if option != nil {
		filePath = option.FilePath
		bufTickInterval = option.BufRefreshTickInterval
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("NewTracker Fail: %v", err))
		return nil
	}
	tracker := &Tracker{
		file:      file,
		buffer:    make([]byte, 0),
		bufTicker: time.NewTicker(time.Duration(bufTickInterval) * time.Millisecond),
		closer:    make(chan bool),
	}
	tracker.start()
	return tracker
}

func (t *Tracker) start() {
	go t.scheduledRefresh()
	log.Println("Tracker's Refresh worker starts working...")
}

func (t *Tracker) Close() {
	t.closer <- true
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

func (t *Tracker) Write(content []byte) {
	t.bufLock.Lock()
	defer t.bufLock.Unlock()
	t.buffer = append(t.buffer, content...)
}

func (t *Tracker) refresh() {
	t.bufLock.Lock()
	defer t.bufLock.Unlock()
	if len(t.buffer) == 0 {
		return
	}
	_, err := t.file.Write(t.buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.buffer = []byte{}
}
