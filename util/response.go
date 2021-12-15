package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type WareResponse struct {
	Code int
	Val  interface{}
	Msg  string
}

func MakeResponse(c *gin.Context, resp *WareResponse) {
	var msg string
	if resp.Msg != "" {
		msg = resp.Msg
	} else {
		msg = ErrCode2Msg[resp.Code]
	}

	respMap := gin.H{
		"code": resp.Code,
		"msg":  msg,
	}
	if resp.Val != nil {
		respMap["data"] = resp.Val
	}
	timeCost := TimeCost(c)
	if timeCost != "" {
		respMap["cost_time"] = timeCost
	}

	c.JSON(http.StatusOK, respMap)
}
