package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	todoBackendAddress := os.Getenv("TODO_BACKEND_ADDRESS")
	if todoBackendAddress == "" {
		todoBackendAddress = "http://todo-app-backend-svc:18083"
		log.Printf("TODO_BACKEND_ADDRESS not set, using default: %s", todoBackendAddress)
	}

	// 1. Get random wikipedia page
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://en.wikipedia.org/wiki/Special:Random", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	// Wikipedia requires a User-Agent header to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch random wikipedia page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code from Wikipedia: %d", resp.StatusCode)
	}

	finalURL := resp.Request.URL.String()
	log.Printf("Get target URL:%s", finalURL);

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	bodyString := string(bodyBytes)

	// Simple title extraction
	titleStart := strings.Index(bodyString, "<title>")
	titleEnd := strings.Index(bodyString, "</title>")
	if titleStart == -1 || titleEnd == -1 {
		log.Fatal("Could not find title in response")
	}
	// <title>Random Page - Wikipedia</title> -> extract content
	pageTitle := bodyString[titleStart+7 : titleEnd]
	// Remove " - Wikipedia" suffix if present for cleaner output
	pageTitle = strings.TrimSuffix(pageTitle, " - Wikipedia")

	log.Printf("Found page: %s (%s)", pageTitle, finalURL)

	// 2. Send POST request to backend
	todoItem := fmt.Sprintf("阅读：%s %s", pageTitle, finalURL)
	jsonData := map[string]string{
		"item": todoItem,
	}
	jsonValue, _ := json.Marshal(jsonData)

	postURL := fmt.Sprintf("%s/todos", todoBackendAddress)
	respPost, err := http.Post(postURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatalf("Failed to post to backend: %v", err)
	}
	defer respPost.Body.Close()

	if respPost.StatusCode >= 200 && respPost.StatusCode < 300 {
		log.Println("Successfully added todo item")
	} else {
		log.Printf("Failed to add todo item, status code: %d", respPost.StatusCode)
	}
}
