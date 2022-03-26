package storage

type KVTable interface {
	Get(key *Key) Value
	Set(key *Key, val Value)
	SetInTime(key *Key, val Value)
	Delete(key *Key)
}
