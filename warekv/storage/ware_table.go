package storage

const (
	DefaultTableNum = 16
)

var (
	wTable *WareTable
)

func init() {
	wTable = &WareTable{}
	wTable.TableList = make([]*Shard, DefaultTableNum)
	wTable.TableNum = DefaultTableNum
	for i := range wTable.TableList {
		wTable.TableList[i] = newShard()
	}
}

// WareTable 总表
type WareTable struct {
	TableList []*Shard
	TableNum  int // 永远保持2的倍数，方便哈希计算
}

func GetWareTable() *WareTable {
	return wTable
}

func (w *WareTable) wHash(key *Key) int {
	hashCode := key.Hashcode()
	// TableNum保持2的倍数，方便hash计算
	// 默认16，16-1=15 --> 二进制表示：1111
	// 通过与运算提高取模效率
	return hashCode & (w.TableNum - 1)
}

func (w *WareTable) Get(key *Key) Value {
	pos := w.wHash(key)
	return w.TableList[pos].Get(key)
}

func (w *WareTable) Set(key *Key, val Value) {
	pos := w.wHash(key)
	w.TableList[pos].Set(key, val)
}

func (w *WareTable) Delete(key *Key) {
	pos := w.wHash(key)
	w.TableList[pos].Delete(key)
}
