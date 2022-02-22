package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"ware-kv/warekv/handler"
	"ware-kv/warekv/manager"
	"ware-kv/warekv/storage"
)

var Server *WareKV

type WareKV struct {
	wTable          *storage.WareTable
	router          *gin.Engine
	subscribeCenter *manager.SubscribeCenter
	//options（布隆开关、日志开关）
	//stats
	//machine info
	// closer
	//tracker(AOF) head|k|v|crc32 mmap
	// camera(RDB)
	// GC
	// 热点采样
}

func Boot(port string) {
	Server = &WareKV{
		wTable:          storage.GetWareTable(),
		router:          gin.Default(),
		subscribeCenter: manager.GetSubscribeCenter(),
	}
	handler.Register(Server.router)
	showFrame()
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
