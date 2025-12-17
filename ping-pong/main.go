package main

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	var counter int64 = 0

	r.GET("/pingpong", func(c *gin.Context) {
		current := atomic.AddInt64(&counter, 1)
		c.String(200, fmt.Sprintf("pong %d", current-1))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
