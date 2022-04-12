package authentication

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		switch method {
		case http.MethodGet:
			gin.BasicAuth(GetAuthCenter().GetReaders())(c)
		case http.MethodPost, http.MethodPut, http.MethodDelete:
			gin.BasicAuth(GetAuthCenter().GetWriters())(c)
		default:
			c.JSON(http.StatusOK, "This HTTP method is not allowed in this version...")
			return
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
