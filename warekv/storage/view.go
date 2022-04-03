package storage

import (
	"encoding/json"
	"github.com/qizong007/ware-kv/warekv/util"
	"log"
)

func kvPairView(key string, value Value) []byte {
	// type (1 byte)
	tipe := byte(value.GetType())

	// key (key len bytes)
	keyBytes := []byte(key)

	// base (base len bytes)
	baseBytes, err := json.Marshal(value.GetBase())
	if err != nil {
		log.Println(key, "base json marsh failed!")
		return []byte{}
	}

	// value (value len bytes)
	valueBytes, err := json.Marshal(value.GetValue())
	if err != nil {
		log.Println(key, "value json marsh failed!")
		return []byte{}
	}

	// key len (4 bytes)
	keyLen := len(keyBytes)

	// base len (4 bytes)
	baseLen := len(baseBytes)

	// value len (4 bytes)
	valueLen := len(valueBytes)

	data := make([]byte, 0, 13+keyLen+baseLen+valueLen)
	data = append(data, tipe)
	data = append(data, util.IntToBytes(keyLen)...)
	data = append(data, keyBytes...)
	data = append(data, util.IntToBytes(baseLen)...)
	data = append(data, baseBytes...)
	data = append(data, util.IntToBytes(valueLen)...)
	data = append(data, valueBytes...)

	return data
}
