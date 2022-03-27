package util

type DSType uint8

const (
	StringDS = iota
	CounterDS
	ObjectDS
	ListDS
	ZListDS
	SetDS
	BitmapDS
	BloomFilterDS
	LockDS
)
