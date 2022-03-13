package tracker

import (
	"encoding/json"
	"log"
	"ware-kv/warekv"
	"ware-kv/warekv/storage"
)

type DeleteCommand struct {
	CommandBase
}

func NewDeleteCommand(key string) *DeleteCommand {
	return &DeleteCommand{CommandBase{key}}
}

func (c *DeleteCommand) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		log.Println("DeleteCommand Json Marshall Fail", err)
		return ""
	}
	return string(data)
}

func (c *DeleteCommand) Execute() {
	warekv.Engine().Delete(storage.MakeKey(c.Key))
}

func (c *DeleteCommand) GetOpType() string {
	return DeleteOp
}
