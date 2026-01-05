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
	ID        int    `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
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

	// 修改表结构，添加 completed 字段
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS todos (id SERIAL PRIMARY KEY, item TEXT, completed BOOLEAN DEFAULT FALSE)")
	if err != nil {
		log.Fatal(err)
	}

	// 1. 初始化 Gin 引擎
	r := gin.Default()
	
	// 存活探针，Gin框架正常运行证明存活
	r.GET("/healthz", func(c *gin.Context) {
        c.Status(http.StatusOK)
    })

	// 新增 /ready 路由用于健康检查
	r.GET("/ready", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "database connection failed"})
			return
		}
		c.Status(http.StatusOK)
	})

	// 2. 设置用于获取todo-app的路由
	r.GET("/todos", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, item, completed FROM todos")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		todoList := []Todo{}
		for rows.Next() {
			var t Todo
			if err := rows.Scan(&t.ID, &t.Item, &t.Completed); err != nil {
				continue
			}
			todoList = append(todoList, t)
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
			rows, err := db.Query("SELECT id, item, completed FROM todos")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()

			todoList := []Todo{}
			for rows.Next() {
				var t Todo
				rows.Scan(&t.ID, &t.Item, &t.Completed)
				todoList = append(todoList, t)
			}

			c.JSON(http.StatusCreated, gin.H{
				"message": "Todo added successfully",
				"todos":   todoList,
			})
		}

	})

	// 4. 设置用于更新todo状态的路由
	r.PUT("/todos/:id", func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec("UPDATE todos SET completed = TRUE WHERE id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	// 6. 启动服务
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
