//go:build generate
// +build generate

package main

import (
	"crypto/rand"
	"fmt"
	"os"
)

func main() {
	// AES-256 使用 32 字节的密钥
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法生成随机密钥: %v\n", err)
		os.Exit(1)
	}

	// fmt.Println(key)
	err := os.WriteFile("encrypt/capture.log.key", key, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法写入 capture.log.key: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("成功生成并保存了新的 AES 密钥到 capture.log.key。")
}
