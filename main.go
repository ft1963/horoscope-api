package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 各エンドポイントにハンドラ関数を登録
	http.HandleFunc("/", handleHello)
	http.HandleFunc("/sun-sign", handleSunSign)
	http.HandleFunc("/htmx-sun-sign", handleHtmxSunSign)
	http.HandleFunc("/ten", handleTen)

	fmt.Printf("Server starting on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}