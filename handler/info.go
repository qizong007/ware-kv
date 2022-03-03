package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ware-kv/warekv/machine"
)

func Info(c *gin.Context) {
	c.JSON(http.StatusOK, machine.GetWareInfo())
}
