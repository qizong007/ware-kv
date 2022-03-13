package anticorrosive

import (
	"ware-kv/warekv"
	"ware-kv/warekv/manager"
	"ware-kv/warekv/storage"
)

func Set(key *storage.Key, newVal storage.Value) {
	warekv.Engine().Set(key, newVal)
	SetNotify(key, newVal)
}

func SetNotify(key *storage.Key, newVal storage.Value) {
	go warekv.Engine().Notify(key.GetKey(), newVal.GetValue(), manager.CallbackSetEvent)
}

func Del(key *storage.Key) {
	warekv.Engine().Delete(key)
	deleteNotify(key)
}

func deleteNotify(key *storage.Key) {
	go warekv.Engine().Notify(key.GetKey(), nil, manager.CallbackDeleteEvent)
}
