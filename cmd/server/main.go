package main

import (
	"fmt"
	"github.com/jiu-u/oai-adapter/clients/gemini"
	"net/http"
	"net/url"
)

func main() {
	proxyURL, err := url.Parse("sock5://127.0.0.1:7890") // 替换为你的代理地址和端口
	if err != nil {
		panic(err)
	}
	geminiClient := gemini.NewClient("https://generativelanguage.googleapis.com",
		"xxxxxxxxxxxxxxxxxxxxxx", proxyURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/chat/completions", HandleChatCompletions(geminiClient))

	handler := corsMiddleware(mux)

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", handler)
}

// CORS中间件
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头部
		w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理预检请求
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r) // 调用下一个处理程序
	})
}
