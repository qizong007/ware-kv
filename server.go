package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"time"
	"ware-kv/camera"
	"ware-kv/handler"
	"ware-kv/tracker"
	"ware-kv/util"
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
	camera  *camera.Camera
}

func Boot(option *WareOption) {
	bootTime := time.Now()
	initOption(option)
	tk := tracker.NewTracker(option.Tracker)
	defer tk.Close()
	cmr := camera.NewCamera(option.Camera)
	defer cmr.Close()
	Server = &WareKV{
		engine:  warekv.New(option.WareEngine),
		router:  gin.Default(),
		tracker: tk,
		camera:  cmr,
	}
	// first load the camera
	cmr.DevelopPhotos()
	// then load the tracker
	tk.LoadTracker()
	// Server.engine start in New()
	defer Server.engine.Close()
	handler.Register(Server.router)
	fmt.Println(" -----------------------------------")
	fmt.Printf("  Ware-KV Loading Cost %s  \n", time.Since(bootTime).String())
	fmt.Println(" -----------------------------------")
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
	color.HiMagenta.Printf("         😎version_%s@qizong007\n", util.WareKVVersion)
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
