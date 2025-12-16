package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
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

func printTimeAndString(targetString string) {
	fmt.Printf("%s: %s\n", time.Now().UTC().Format(time.RFC3339Nano), targetString)
}

func main() {
	id, err := uuidV4()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to generate id:", err)
		os.Exit(1)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// 立即进行一次输出
	printTimeAndString(id)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			printTimeAndString(id)
		case <-sigs:
			// 优雅退出
			return
		}
	}
}
