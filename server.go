package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"github.com/qizong007/ware-kv/authentication"
	"github.com/qizong007/ware-kv/camera"
	"github.com/qizong007/ware-kv/handler"
	"github.com/qizong007/ware-kv/tracker"
	"github.com/qizong007/ware-kv/util"
	"github.com/qizong007/ware-kv/warekv"
	"log"
	"time"
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
	gin.SetMode(gin.ReleaseMode)
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
	needAuth := registerAuth(option.Auth)
	handler.Register(Server.router, needAuth)
	showFrame()
	port := defaultPort
	if option != nil {
		port = option.Port
	}
	fmt.Println(" -----------------------------------------------------------")
	fmt.Printf("  Ware-KV Loading Cost %s, listening on port:%s \n", time.Since(bootTime).String(), port)
	fmt.Println(" -----------------------------------------------------------")
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

func registerAuth(option *authentication.WareAuthOption) bool {
	if option == nil || option.Open == false {
		return false
	}
	authCenter := authentication.NewAuthCenter(&authentication.AuthCenterOption{
		Username: option.Root.Username,
		Password: option.Root.Password,
	})
	if option.Others != nil && len(option.Others) > 0 {
		for _, user := range option.Others {
			if err := authCenter.Register(&authentication.AuthRegisterOption{
				ParentUser: option.Root.Username,
				Username:   user.Username,
				Password:   user.Password,
				Role:       authentication.GetAuthRoleIntFromStr(user.Role),
			}); err != nil {
				log.Println(user.Username, "Register Fail:", err)
			}
		}
	}
	return true
}
