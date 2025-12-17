package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	filePath = "../../files/random_string.txt"
	interval = 5 * time.Second
)

func generateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	// Generate a random string at startup
	randomStr := generateRandomString(16)
	fmt.Printf("Generated random string: %s\n", randomStr)
	fmt.Printf("Writing to %s every %s...\n", filePath, interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			continue
		}

		timestamp := time.Now().Format(time.RFC3339)
		content := fmt.Sprintf("[%s] %s\n", timestamp, randomStr)

		if _, err := f.WriteString(content); err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
		} else {
			fmt.Println("Wrote to file successfully.")
		}
		f.Close()
	}
}
