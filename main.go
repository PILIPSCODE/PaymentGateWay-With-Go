package main

import (
	"backend/controlers/productcontrollers"

	"github.com/gin-gonic/gin"
	"subosito.com/go/gotenv"
)

func main() {
	gotenv.Load()
	r := gin.Default()

	r.Use(corsMiddleware())
	r.POST("/api/product", productcontrollers.Initial)

	r.Run()
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
