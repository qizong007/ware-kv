package storage

import (
	"sync"
)

// TableUnit 表单元
type TableUnit struct {
	Table map[string]Value
	TLock sync.RWMutex
}

func newTableUnit() *TableUnit {
	table := &TableUnit{}
	table.Table = make(map[string]Value)
	return table
}

func (t *TableUnit) Get(key *Key) Value {
	t.TLock.RLock()
	defer t.TLock.RUnlock()
	return t.Table[key.Val]
}

func (t *TableUnit) Set(key *Key, val Value) {
	t.TLock.Lock()
	defer t.TLock.Unlock()
	t.Table[key.Val] = val
}

func (t *TableUnit) Delete(key *Key) {
	t.TLock.Lock()
	defer t.TLock.Unlock()
	t.Table[key.Val].DeleteValue()
}
