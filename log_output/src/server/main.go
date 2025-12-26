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
	envMessage := os.Getenv("MESSAGE")
	if envMessage == "" {
		envMessage = "hello world"
	}

	configDirPath := os.Getenv("CONFIG_DIR_PATH")
	if configDirPath == "" {
		configDirPath = "./config"
	}

	randomStringFilePath := os.Getenv("RANDOM_STRING_FILE_PATH")
	if randomStringFilePath == "" {
		randomStringFilePath = "./random_string.txt"
	}

	pingPongServerAddress := os.Getenv("PING_PONG_SERVER_ADDRESS")
	if pingPongServerAddress == "" {
		pingPongServerAddress = "http://ping-pong-svc:18081"
	}

	// 创建Web服务器
	r := gin.Default()

	// 创建请求处理器
	r.GET("/", func(c *gin.Context) {
		informationFile, err := os.ReadFile(configDirPath + "/information.txt")
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			c.String(http.StatusInternalServerError, "Error opening file: %v", err)
			return
		}

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

		c.String(http.StatusOK,
			"file content: %s\nenv variable: %s\n%s\nPing / Pongs: %s", informationFile, envMessage, randomString, pingCount)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
