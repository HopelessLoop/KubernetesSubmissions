package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	randomStringFilePath := os.Getenv("RANDOM_STRING_FILE_PATH")
	if randomStringFilePath == "" {
		randomStringFilePath = "./random_string.txt"
	}

	// pingPongFilePath := os.Getenv("PING_PONG_FILE_PATH")
	// if pingPongFilePath == "" {
	// 	pingPongFilePath = "./ping-pong.txt"
	// }

	pingPongServerAddress := os.Getenv("PING_PONG_SERVER_ADDRESS")
	if pingPongServerAddress == "" {
		pingPongServerAddress = "http://ping-pong-svc:18081"
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

		// 尝试从 HTTP 服务获取
		var pingCount string
		resp, err := http.Get(pingPongServerAddress + "/pings")
		if err == nil {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			var result map[string]int64
			if json.Unmarshal(body, &result) == nil {
				pingCount = fmt.Sprintf("%d", result["pings"])
			}
		}

		// 如果 HTTP 获取失败，尝试从文件读取
		// if pingCount == "" {
		// 	pingPongs, err := os.ReadFile(pingPongFilePath)
		// 	if err != nil {
		// 		// 如果文件不存在，可能还没有ping-pong请求，默认为0或空
		// 		pingCount = "0"
		// 	} else {
		// 		pingCount = string(pingPongs)
		// 	}
		// }

		c.String(http.StatusOK, "%s\nPing / Pongs: %s", randomString, pingCount)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
