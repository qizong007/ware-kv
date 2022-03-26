package storage

import "ware-kv/warekv/closer"

type KVTable interface {
	Get(key *Key) Value
	Set(key *Key, val Value)
	SetInTime(key *Key, val Value)
	Delete(key *Key)
	closer.Closer // for resource collection
	KeyNum() int
	Type() string
}

var GlobalTable KVTable
