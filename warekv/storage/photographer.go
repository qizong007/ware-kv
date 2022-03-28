package storage

type Photographer interface {
	View() []byte
	GetFlag() Flag
}

type Flag int

const (
	TableFlag = iota
	SubscribeCenterFlag
)
