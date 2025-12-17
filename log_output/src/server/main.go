package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	filePath = "../../files/random_string.txt"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)

	// 创建Web服务器
	r := gin.Default()

	// 创建请求处理器
	r.GET("/", func(c *gin.Context) {
		randomString, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			c.String(http.StatusInternalServerError, "Error opening file: %v", err)
		} else {
			c.String(http.StatusOK, "%s", randomString)
		}
	})

	r.Run(":" + port)
}
