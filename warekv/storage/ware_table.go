package storage

import (
	"fmt"
	"log"
	"time"
	"ware-kv/warekv/util"
)

const (
	defaultShardNum          = 16
	defaultWriteQueueCap     = 256 // 默认写请求缓存容量
	defaultWriteTickInterval = 100 * time.Millisecond
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

func NewWareTable(shardOption *ShardOption, gcOption *WareGCOption) *WareTable {
	wTable = &WareTable{}
	shardNum := defaultShardNum
	writeQueueCap := defaultWriteQueueCap
	writeTickInterval := defaultWriteTickInterval
	if shardOption != nil {
		shardNum = int(util.Nearest2Power(shardOption.Num))
		writeQueueCap = int(shardOption.WriteQueueCap)
		writeTickInterval = time.Millisecond * time.Duration(shardOption.WriteTickInterval)
	}
	wTable.TableList = make([]*Shard, shardNum)
	wTable.TableNum = shardNum
	for i := range wTable.TableList {
		wTable.TableList[i] = newShard(writeQueueCap, writeTickInterval, gcOption)
	}
	return wTable
}

func (w *WareTable) Start() {
	for _, shard := range w.TableList {
		shard.Start()
	}
	log.Println("WareTable's Write worker and GC worker start working...")
}

func (w *WareTable) Close() {
	for _, shard := range w.TableList {
		shard.Close()
	}
	log.Println("WareTable's Write worker and GC worker stop working...")
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

// Delete 标记删除
func (w *WareTable) Delete(key *Key) {
	pos := w.wHash(key)
	w.TableList[pos].Delete(key)
}
