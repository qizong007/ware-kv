package tracker

import "fmt"

type TrackOp string

const (
	CreateOp = "c"
	ModifyOp = "m"
	DeleteOp = "d"
)

type Command interface {
	fmt.Stringer
	Execute()
	GetOpType() string
}

type CommandBase struct {
	Key string `json:"k"`
}
