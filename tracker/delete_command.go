package tracker

import (
	"encoding/json"
	"github.com/qizong007/ware-kv/warekv"
	"github.com/qizong007/ware-kv/warekv/storage"
	"log"
)

type DeleteCommand struct {
	Key string `json:"k"`
}

func NewDeleteCommand(key string) *DeleteCommand {
	return &DeleteCommand{key}
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
	warekv.Engine().DeleteInTime(storage.MakeKey(c.Key))
}

func (c *DeleteCommand) GetOpType() string {
	return DeleteOp
}
