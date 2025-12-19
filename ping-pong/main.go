package main

import (
	"fmt"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

func loadCounter(filePath string) int64 {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0
	}

	val, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0
	}

	return val
}

func saveCounter(filePath string, val int64) {
	os.WriteFile(filePath, []byte(fmt.Sprintf("%d", val)), 0644)
}

func main() {
	pingPongFilePath := os.Getenv("PING_PONG_FILE_PATH")
	if pingPongFilePath == "" {
		pingPongFilePath = "./ping-pong.txt"
	}

	var counter int64 = loadCounter(pingPongFilePath)

	r := gin.Default()
	r.GET("/pingpong", func(c *gin.Context) {
		current := atomic.AddInt64(&counter, 1)
		saveCounter(pingPongFilePath, current)
		c.String(200, fmt.Sprintf("pong %d", current))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
