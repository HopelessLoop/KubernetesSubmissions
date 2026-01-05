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

	todoBackendAddress := os.Getenv("TODO_BACKEND_ADDRESS")
	if todoBackendAddress == "" {
		todoBackendAddress = "http://todo-app-backend-svc:18083"
		log.Printf("TODO_BACKEND_ADDRESS not set, using default: %s", todoBackendAddress)
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

	// 5. 设置用于获取图像路由
	r.GET("/image", func(c *gin.Context) {
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			c.String(http.StatusServiceUnavailable, "Image not ready yet")
			return
		}
		c.File(imagePath)
	})

	// 5. 设置用于获取todo-app的路由
	r.GET("/todo-app", func(c *gin.Context) {
		html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Todo APP</title>
    <script>
        const backendUrl = "%s";

        async function fetchTodos() {
            try {
                const response = await fetch(backendUrl + "/todos");
                if (!response.ok) throw new Error('Failed to fetch');
                const data = await response.json();
                
                const todoList = document.getElementById('todo-list');
                const doneList = document.getElementById('done-list');
                todoList.innerHTML = '';
                doneList.innerHTML = '';

                if (data.todos) {
                    data.todos.forEach(todo => {
                        const li = document.createElement('li');
                        li.textContent = todo.item;
                        
                        if (todo.completed) {
                            doneList.appendChild(li);
                        } else {
                            const button = document.createElement('button');
                            button.textContent = 'Mark as done';
                            button.onclick = () => markAsDone(todo.id);
                            button.style.marginLeft = '10px';
                            li.appendChild(button);
                            todoList.appendChild(li);
                        }
                    });
                }
            } catch (error) {
                console.error('Error fetching todos:', error);
            }
        }

        async function createTodo() {
            const input = document.getElementById('todo-input');
            const item = input.value;
            if (!item) return;

            try {
                const response = await fetch(backendUrl + "/todos", {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ item: item }),
                });
                if (response.ok) {
                    input.value = '';
                    fetchTodos();
                } else {
                    console.error('Failed to create todo');
                }
            } catch (error) {
                console.error('Error creating todo:', error);
            }
        }

        async function markAsDone(id) {
            try {
                const response = await fetch(backendUrl + "/todos/" + id, {
                    method: 'PUT'
                });
                if (response.ok) {
                    fetchTodos();
                } else {
                    console.error('Failed to mark as done');
                }
            } catch (error) {
                console.error('Error marking as done:', error);
            }
        }

        window.onload = fetchTodos;
    </script>
</head>
<body>
    <h1>Todo APP</h1>
    
    <!-- 图片区域 -->
    <img src="/image" alt="Daily Image" style="max-width: 500px; display: block; margin-bottom: 20px;">
    
    <!-- 输入区域 -->
    <div>
        <input id="todo-input" type="text" maxlength="140" placeholder="Enter a todo (max 140 chars)">
        <button type="button" onclick="createTodo()">Create todo</button>
    </div>

    <!-- 待办事项列表 -->
    <h2>To Do</h2>
    <ul id="todo-list">
        <!-- Active items will be loaded here -->
    </ul>

    <h2>Done</h2>
    <ul id="done-list">
        <!-- Completed items will be loaded here -->
    </ul>
</body>
</html>`, todoBackendAddress)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	// 6. 启动服务
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
