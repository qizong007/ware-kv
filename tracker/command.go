package tracker

import "fmt"

type TrackOp string

const (
	CreateOp    = "c"
	ModifyOp    = "m"
	DeleteOp    = "d"
	SubscribeOp = "s"
)

type Command interface {
	fmt.Stringer
	Execute()
	GetOpType() string
}
