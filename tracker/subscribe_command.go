package tracker

import (
	"encoding/json"
	"log"
	"ware-kv/warekv"
	"ware-kv/warekv/manager"
)

type SubCommand struct {
	Manifest *manager.SubscribeManifest `json:"mf"`
}

func NewSubCommand(manifest *manager.SubscribeManifest) *SubCommand {
	return &SubCommand{manifest}
}

func (c *SubCommand) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		log.Println("SubCommand Json Marshall Fail", err)
		return ""
	}
	return string(data)
}

func (c *SubCommand) Execute() {
	warekv.Engine().Subscribe(c.Manifest)
}

func (c *SubCommand) GetOpType() string {
	return SubscribeOp
}
