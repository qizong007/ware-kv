package log

import (
	"encoding/json"
	"fmt"
	"log"
)

type TrackOp string

const (
	CreateOp = "c"
	ModifyOp = "m"
	DeleteOp = "d"
)

type Command interface {
	String() string
}

type CommandBase struct {
	Op  TrackOp
	Key string
}

func (c *CommandBase) String() string {
	return string(c.Op) + " " + c.Key
}

type CreateCommand struct {
	CommandBase
	Val        interface{}
	CreateTime int64
	ExpireTime int64
}

func GenCreateCommand(key string, val interface{}, createTime int64, expireTime int64) *CreateCommand {
	return &CreateCommand{
		CommandBase: CommandBase{CreateOp,key},
		Val:         val,
		CreateTime:  createTime,
		ExpireTime:  expireTime,
	}
}

func (c *CreateCommand) String() string {
	valStr, err := json.Marshal(c.Val)
	if err != nil {
		log.Println("CreateCommand Json Marshall Fail", err)
		return ""
	}
	str := fmt.Sprintf("%s %s %d %d", c.CommandBase.String(), string(valStr), c.CreateTime, c.ExpireTime)
	return str
}

type DeleteCommand struct {
	CommandBase
}

func GenDeleteCommand(key string) *DeleteCommand {
	return &DeleteCommand{CommandBase{DeleteOp,key}}
}

func (c *DeleteCommand) String() string {
	return c.CommandBase.String()
}
