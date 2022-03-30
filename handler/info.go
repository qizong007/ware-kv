package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qizong007/ware-kv/warekv/machine"
	"net/http"
)

func Info(c *gin.Context) {
	c.JSON(http.StatusOK, machine.GetWareInfo())
}
