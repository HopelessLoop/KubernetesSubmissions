package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 获取图片存储路径
	imagePath := os.Getenv("IMAGE_PATH")
	if imagePath == "" {
		imagePath = "./image.jpg"
		log.Printf("IMAGE_PATH not set, using default: %s", imagePath)
	}

	// 2. 定义更新图片的函数
	updateImage := func() {
		log.Println("Starting image update...")
		resp, err := http.Get("https://picsum.photos/1200")
		if err != nil {
			log.Printf("Failed to download image: %v", err)
			return
		}
		defer resp.Body.Close()

		file, err := os.Create(imagePath)
		if err != nil {
			log.Printf("Failed to create file: %v", err)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			log.Printf("Failed to save image: %v", err)
			return
		}
		log.Println("Image updated successfully.")
	}

	// 3. 启动定时任务 (每小时更新一次)
	go func() {
		// 初始化时立即更新一次
		updateImage()
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			updateImage()
		}
	}()

	// 4. 初始化 Gin 引擎
	r := gin.Default()

	// 5. 设置路由
	r.GET("/image", func(c *gin.Context) {
		// 检查文件是否存在
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			c.String(http.StatusServiceUnavailable, "Image not ready yet")
			return
		}
		// 返回图片文件
		c.File(imagePath)
	})

	// 6. 启动服务
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
