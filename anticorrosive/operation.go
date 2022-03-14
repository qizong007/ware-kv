package anticorrosive

import (
	"log"
	"ware-kv/util"
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

func IsKVEffective(val storage.Value) (bool, int) {
	if val == nil {
		log.Println("key is not existed")
		return false, util.KeyNotExisted
	}
	if !val.IsAlive() {
		log.Println("key is dead")
		return false, util.KeyHasDeleted
	}
	if val.IsExpired() {
		log.Println("key has been expired")
		return false, util.KeyHasExpired
	}
	return true, 0
}
