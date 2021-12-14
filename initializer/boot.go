package initializer

import (
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
)

func Boot() *gin.Engine {
	engine := gin.Default()
	Register(engine)
	showFrame()
	return engine
}

func showFrame() {
	color.HiYellow.Println("                                   __")
	color.HiYellow.Println("   _      ______ _________        / /___   __")
	color.HiGreen.Println("  | | /| / / __ `/ ___/ _ \\______/ //_/ | / /")
	color.HiCyan.Println("  | |/ |/ / /_/ / /  /  __/_____/ ,<  | |/ /")
	color.HiBlue.Print("  |__/|__/\\__,_/_/   \\___/     /_/|_| |___/")
	color.HiMagenta.Println("         😎version_0.0.1@qizong007")
}
