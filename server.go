package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"ware-kv/handler"
	"ware-kv/tracker"
	"ware-kv/warekv"
)

const (
	defaultPort = "7777"
)

var Server *WareKV

type WareKV struct {
	engine  *warekv.WareEngine
	router  *gin.Engine
	tracker *tracker.Tracker
}

func Boot(option *WareOption) {
	initOption(option)
	tk := tracker.NewTracker(option.Tracker)
	Server = &WareKV{
		engine:  warekv.New(option.WareEngine),
		router:  gin.Default(),
		tracker: tk,
	}
	defer tk.Close()
	tk.LoadTracker()
	// Server.engine start in New()
	defer Server.engine.Close()
	handler.Register(Server.router)
	showFrame()
	port := defaultPort
	if option != nil {
		port = option.Port
	}
	Server.router.Run(fmt.Sprintf(":%s", port))
}

func showFrame() {
	color.HiYellow.Println("                                   __")
	color.HiYellow.Println("   _      ______ _________        / /___   __")
	color.HiGreen.Println("  | | /| / / __ `/ ___/ _ \\______/ //_/ | / /")
	color.HiCyan.Println("  | |/ |/ / /_/ / /  /  __/_____/ ,<  | |/ /")
	color.HiBlue.Print("  |__/|__/\\__,_/_/   \\___/     /_/|_| |___/")
	color.HiMagenta.Println("         😎version_0.0.1@qizong007")
}

func initOption(option *WareOption) {
	engineOption := DefaultOption()
	if option.WareEngine != nil {
		if option.WareEngine.GC == nil {
			option.WareEngine.GC = engineOption.WareEngine.GC
		}
		if option.WareEngine.Shard == nil {
			option.WareEngine.Shard = engineOption.WareEngine.Shard
		}
		if option.WareEngine.Subscriber == nil {
			option.WareEngine.Subscriber = engineOption.WareEngine.Subscriber
		}
		if option.WareEngine.MachineInfo == nil {
			option.WareEngine.MachineInfo = engineOption.WareEngine.MachineInfo
		}
	}
}
