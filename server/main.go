package main

import (
	"github.com/gin-gonic/gin"
	"rsc.io/quote"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "hello world")
	})

	r.GET("/quote", func(c *gin.Context) {
		c.JSON(200, quote.Go())
	})

	r.Run(":3030")
}
