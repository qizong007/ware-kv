package util

import (
	"github.com/gin-gonic/gin"
	"time"
)

const (
	TimeTag = "t"
)

func TimeKeeping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(TimeTag, time.Now())
		c.Next()
	}
}

func TimeCost(c *gin.Context) string {
	if startTime, ok := c.Get(TimeTag); !ok {
		return ""
	} else {
		return time.Since(startTime.(time.Time)).String()
	}
}
