package cache

import "ware-kv/warekv/storage"

type Cache interface {
	Get(*storage.Key) (storage.Value, bool)
	Set(*storage.Key, storage.Value)
	Delete(*storage.Key)
}

type entry struct {
	key   *storage.Key
	value storage.Value
}
