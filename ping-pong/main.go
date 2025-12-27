package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func initDB(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS pings (id SERIAL PRIMARY KEY, count INT)")
	if err != nil {
		log.Fatal(err)
	}

	var count int
	err = db.QueryRow("SELECT count(*) FROM pings").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		_, err = db.Exec("INSERT INTO pings (count) VALUES (0)")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	postgresAddress := os.Getenv("POSTGRES_ADDRESS")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")

	if postgresAddress == "" {
		log.Fatal("POSTGRES_ADDRESS environment variable required")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s/postgres?sslmode=disable", postgresUser, postgresPassword, postgresAddress)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initDB(db)

	r := gin.Default()
	r.GET("/pingpong", func(c *gin.Context) {
		var current int
		err := db.QueryRow("UPDATE pings SET count = count + 1 RETURNING count").Scan(&current)
		if err != nil {
			c.String(500, "Database error")
			return
		}
		c.String(200, fmt.Sprintf("pong %d", current))
	})

	r.GET("/pings", func(c *gin.Context) {
		var val int
		err := db.QueryRow("SELECT count FROM pings LIMIT 1").Scan(&val)
		if err != nil {
			c.String(500, "Database error")
			return
		}
		c.JSON(200, gin.H{
			"pings": val,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
