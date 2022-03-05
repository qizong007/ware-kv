package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"ware-kv/handler"
	"ware-kv/warekv"
)

const (
	defaultPort = "7777"
)

var Server *WareKV

type WareKV struct {
	engine *warekv.WareEngine
	router *gin.Engine
	// options（布隆开关、日志开关）
	// tracker(AOF) head|k|v|crc32 mmap
}

func Boot(option *WareOption) {
	Server = &WareKV{
		engine: warekv.New(option.WareEngine),
		router: gin.Default(),
	}
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
