package ds

import "ware-kv/warekv/util"

// ZList 有序列表
type ZList struct {
	Base
	skipList *util.SkipList
}
