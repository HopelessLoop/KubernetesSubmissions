package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// 生成 UUID v4
func uuidV4() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

func main() {
	id, err := uuidV4()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to generate id:", err)
		os.Exit(1)
	}

	// 启动后台协程每5秒输出一次
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		// 立即进行一次输出
		fmt.Printf("%s: %s\n", time.Now().UTC().Format(time.RFC3339Nano), id)

		for range ticker.C {
			fmt.Printf("%s: %s\n", time.Now().UTC().Format(time.RFC3339Nano), id)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)

	// 创建Web服务器
	r := gin.Default()

	// 创建请求处理器
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "%s: %s", time.Now().UTC().Format(time.RFC3339Nano), id)
	})

	r.Run(":" + port)
}
