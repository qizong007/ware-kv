package cache

import "ware-kv/warekv/storage"

type Cache interface {
	Get(*storage.Key) storage.Value
	Set(*storage.Key, storage.Value)
	Delete(*storage.Key)
}
