package handler

import (
	"github.com/gin-gonic/gin"
)

func GetStr(c *gin.Context) {
	paramKey := c.Param("key")
	if paramKey == "" {

	}
	//key := storage.MakeKey(paramKey)
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
