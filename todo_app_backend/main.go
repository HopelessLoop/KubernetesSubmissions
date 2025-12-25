package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// 定义接收数据的结构体
type Todo struct {
	Item string `json:"item"`
}

func main() {
	// 初始化为空切片
	todoList := []string{}

	// 1. 初始化 Gin 引擎
	r := gin.Default()

	// 2. 设置用于获取todo-app的路由
	r.GET("/todos", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"todos": todoList,
		})
	})

	// 3. 设置用于添加todo-app的路由
	r.POST("/todos", func(c *gin.Context) {
		var newTodo Todo
		if err := c.ShouldBindJSON(&newTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if newTodo.Item == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Format error",
			})

		} else {
			todoList = append(todoList, newTodo.Item)
			c.JSON(http.StatusCreated, gin.H{
				"message": "Todo added successfully",
				"todos":   todoList,
			})
		}

	})

	// 6. 启动服务
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
