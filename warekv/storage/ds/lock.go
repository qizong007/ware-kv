package ds

import (
	"fmt"
	"github.com/qizong007/ware-kv/warekv/util"
	"sync"
	"time"
)

type Lock struct {
	Base
	State     bool
	TimeLimit int64 // if hit TimeLimit, the lock will release
	Guid      string
	rw        sync.RWMutex
}

type LockView struct {
	IsLocked bool  `json:"is_locked"`
	TimeLeft int64 `json:"time_left"`
}

func (l *Lock) GetValue() interface{} {
	l.rw.RLock()
	defer l.rw.RUnlock()

	var t int64 = 0
	if l.State {
		t = l.UpdateTime + l.TimeLimit - time.Now().Unix()
		if t < 0 {
			t = 0
		}
	}

	return &LockView{
		IsLocked: l.State,
		TimeLeft: t,
	}
}

func MakeLock() *Lock {
	return &Lock{
		Base: *NewBase(util.LockDS),
	}
}

func (l *Lock) Lock(t int64, guid string) error {
	if l.isLocked() {
		return fmt.Errorf("LockFailed")
	}
	l.rw.Lock()
	defer l.rw.Unlock()
	l.Update()
	l.State = true
	l.TimeLimit = t
	l.Guid = guid
	// time over, release lock
	time.AfterFunc(time.Second*time.Duration(t), func() {
		_ = l.Unlock(l.Guid)
	})
	return nil
}

func (l *Lock) Unlock(guid string) error {
	if !l.isLocked() {
		return fmt.Errorf("UnlockFailed")
	}
	if !l.checkLockKeeper(guid) {
		return fmt.Errorf("NotTheCorrectLockKeeper")
	}
	l.rw.Lock()
	defer l.rw.Unlock()
	l.Update()
	l.State = false
	l.TimeLimit = 0
	return nil
}

func (l *Lock) isLocked() bool {
	l.rw.RLock()
	defer l.rw.RUnlock()
	return l.State
}

func (l *Lock) checkLockKeeper(guid string) bool {
	l.rw.RLock()
	defer l.rw.RUnlock()
	return l.Guid == guid
}
