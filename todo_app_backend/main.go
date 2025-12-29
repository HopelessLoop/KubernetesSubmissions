package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// 定义接收数据的结构体
type Todo struct {
	Item string `json:"item"`
}

func main() {
	postgresAddress := os.Getenv("POSTGRES_ADDRESS")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")

	connStr := fmt.Sprintf("postgres://%s:%s@%s/postgres?sslmode=disable", postgresUser, postgresPassword, postgresAddress)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS todos (id SERIAL PRIMARY KEY, item TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// 1. 初始化 Gin 引擎
	r := gin.Default()

	// 2. 设置用于获取todo-app的路由
	r.GET("/todos", func(c *gin.Context) {
		rows, err := db.Query("SELECT item FROM todos")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		todoList := []string{}
		for rows.Next() {
			var item string
			if err := rows.Scan(&item); err != nil {
				continue
			}
			todoList = append(todoList, item)
		}

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

		} else if len(newTodo.Item) > 140 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Todo item too long (max 140 characters)",
			})
		} else {
			_, err := db.Exec("INSERT INTO todos (item) VALUES ($1)", newTodo.Item)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// 获取最新列表
			rows, err := db.Query("SELECT item FROM todos")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()

			todoList := []string{}
			for rows.Next() {
				var item string
				rows.Scan(&item)
				todoList = append(todoList, item)
			}

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
