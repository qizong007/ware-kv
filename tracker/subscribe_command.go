package tracker

import (
	"encoding/json"
	"github.com/qizong007/ware-kv/warekv"
	"github.com/qizong007/ware-kv/warekv/manager"
	"log"
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
