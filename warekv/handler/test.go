package handler

import (
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func Err500(c *gin.Context) {
	c.JSON(500, gin.H{
		"message": "err",
	})
}
