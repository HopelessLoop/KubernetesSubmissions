package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	randomStringFilePath := os.Getenv("RANDOM_STRING_FILE_PATH")
	if randomStringFilePath == "" {
		randomStringFilePath = "./random_string.txt"
	}

	pingPongFilePath := os.Getenv("PING_PONG_FILE_PATH")
	if pingPongFilePath == "" {
		pingPongFilePath = "./ping-pong.txt"
	}

	// 创建Web服务器
	r := gin.Default()

	// 创建请求处理器
	r.GET("/", func(c *gin.Context) {
		randomString, err := os.ReadFile(randomStringFilePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			c.String(http.StatusInternalServerError, "Error opening file: %v", err)
			return
		}

		pingPongs, err := os.ReadFile(pingPongFilePath)
		if err != nil {
			// 如果文件不存在，可能还没有ping-pong请求，默认为0或空
			pingPongs = []byte("0")
		}

		c.String(http.StatusOK, "%s\nPing / Pongs: %s", randomString, pingPongs)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
