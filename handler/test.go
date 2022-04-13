package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func TestPost(c *gin.Context) {
	params := map[string]interface{}{}
	_ = c.BindJSON(&params)
	fmt.Println(params)
	c.JSON(200, params)
}

func Err500(c *gin.Context) {
	c.JSON(500, gin.H{
		"message": "err",
	})
}
