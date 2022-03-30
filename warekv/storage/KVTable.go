package storage

import "github.com/qizong007/ware-kv/warekv/closer"

type KVTable interface {
	Get(key *Key) Value
	Set(key *Key, val Value)
	SetInTime(key *Key, val Value)
	Delete(key *Key)
	KeyNum() int
	Type() string
	closer.Closer // for resource collection
	Photographer
}

var GlobalTable KVTable
