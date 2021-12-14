package storage

import "sync"

const (
	DefaultTableNum = 16
)

// WareTable 总表
type WareTable struct {
	TableList []*TableUnit
	TableNum  int64
}

// TableUnit 表单元
type TableUnit struct {
	Table map[Key]Value
	TLock sync.RWMutex
}

func newTableUnit() *TableUnit {
	table := &TableUnit{}
	table.Table = make(map[Key]Value)
	return table
}

func NewWareTable() *WareTable {
	wt := &WareTable{}
	wt.TableList = make([]*TableUnit, DefaultTableNum)
	for i := range wt.TableList {
		wt.TableList[i] = newTableUnit()
	}
	return wt
}

func (w WareTable) Get(key *Key) Value {
	return nil
}

func (w WareTable) wHash() {

}
