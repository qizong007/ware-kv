package util

import (
	"bytes"
	"encoding/binary"
)

func IntToBytes(n int) []byte {
	x := int32(n) // 4 bytes
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

func Int64ToBytes(n int64) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, n)
	return bytesBuffer.Bytes()
}

func BytesToInt64(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int64
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}
