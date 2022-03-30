package storage

import (
	"github.com/qizong007/ware-kv/warekv/util"
	"log"
	"time"
)

const (
	defaultShardNum          = 16
	shardNumMin              = 8
	shardNumMax              = 16 * 1024
	defaultWriteQueueCap     = 256
	writeQueueCapMin         = 128
	writeQueueCapMax         = 64 * 1024 * 1024
	defaultWriteTickInterval = 100
	writeTickIntervalMin     = 50
	writeTickIntervalMax     = 1000
)

var (
	wTable *WareTable
)

// WareTable 总表
type WareTable struct {
	TableList []*Shard
	TableNum  int // 永远保持2的倍数，方便哈希计算
}

type ShardOption struct {
	Num               uint `yaml:"Num"`
	WriteQueueCap     uint `yaml:"WriteQueueCap"`
	WriteTickInterval uint `yaml:"WriteTickInterval"`
}

func DefaultShardOption() *ShardOption {
	return &ShardOption{
		Num:               defaultShardNum,
		WriteQueueCap:     defaultWriteQueueCap,
		WriteTickInterval: defaultWriteTickInterval,
	}
}

func NewWareTable(shardOption *ShardOption, gcOption *WareGCOption) *WareTable {
	wTable = &WareTable{}
	shardNum := defaultShardNum
	writeQueueCap := defaultWriteQueueCap
	writeTickInterval := time.Millisecond * time.Duration(defaultWriteTickInterval)
	if shardOption != nil {
		shardNum = util.SetIfHitLimit(int(util.Next2Power(shardOption.Num)), shardNumMin, shardNumMax)
		writeQueueCap = util.SetIfHitLimit(int(shardOption.WriteQueueCap), writeQueueCapMin, writeQueueCapMax)
		tickInterval := util.SetIfHitLimit(int(shardOption.WriteTickInterval), writeTickIntervalMin, writeTickIntervalMax)
		writeTickInterval = time.Millisecond * time.Duration(tickInterval)
	}
	wTable.TableList = make([]*Shard, shardNum)
	wTable.TableNum = shardNum
	for i := range wTable.TableList {
		wTable.TableList[i] = newShard(writeQueueCap, writeTickInterval, gcOption)
	}
	wTable.start()
	return wTable
}

func (w *WareTable) start() {
	for _, shard := range w.TableList {
		shard.start()
	}
	log.Println("WareTable's Write worker and GC worker start working...")
}

func (w *WareTable) Close() {
	for _, shard := range w.TableList {
		shard.Close()
	}
	log.Println("WareTable's Write worker and GC worker stop working...")
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

func (w *WareTable) SetInTime(key *Key, val Value) {
	pos := w.wHash(key)
	w.TableList[pos].SetInTime(key, val)
}

// Delete Mark-Sweep
func (w *WareTable) Delete(key *Key) {
	pos := w.wHash(key)
	w.TableList[pos].Delete(key)
}

func (w *WareTable) KeyNum() int {
	sum := 0
	for i := 0; i < w.TableNum; i++ {
		sum += w.TableList[i].Count()
	}
	return sum
}

func (w *WareTable) Type() string {
	return "WareTable"
}

func (w *WareTable) View() []byte {
	data := make([]byte, 0)
	// table flag
	data = append(data, uint8(TableFlag))
	// keys num
	data = append(data, util.IntToBytes(w.KeyNum())...)
	// kv pairs
	for _, shard := range w.TableList {
		data = append(data, shard.view()...)
	}
	return data
}

func (w *WareTable) GetFlag() Flag {
	return TableFlag
}
