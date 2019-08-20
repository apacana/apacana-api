package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func AddAllowOriginMW(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin")
	fmt.Println(method, origin)
	if origin == "" {
		return
	}

	accessHeaders := c.Request.Header.Get("Access-Control-Request-Headers")
	if accessHeaders != "" {
		c.Header("Access-Control-Allow-Headers", accessHeaders)
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, HEAD")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Max-Age", "86400")

	//放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.AbortWithStatus(200)
	}
}
