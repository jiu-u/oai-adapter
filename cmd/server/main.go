package main

import (
	"flag"
	"fmt"
	oaiadapter "github.com/jiu-u/oai-adapter"
	"github.com/joho/godotenv"
	"net/http"
	"net/url"
	"os"
)

func GetClient() (oaiadapter.Adapter, error) {
	proxyENV := os.Getenv("PROXY")
	clientType := os.Getenv("OAI_TYPE")
	clientURL := os.Getenv("OAI_URL")
	clientKey := os.Getenv("OAI_KEY")

	// 公共配置
	config := &oaiadapter.AdapterConfig{
		AdapterType:  oaiadapter.AdapterType(clientType),
		ApiKey:       clientKey,
		EndPoint:     clientURL,
		ManualModels: nil,
		ProxyURL:     nil,
	}

	if proxyENV != "" {
		proxyURL, err := url.Parse(proxyENV)
		if err != nil {
			return nil, err
		}
		config.ProxyURL = proxyURL
	}

	client := oaiadapter.NewAdapter(config)
	return client, nil
}

func main() {
	devMode := flag.Bool("dev", false, "dev mode")
	flag.Parse()
	if *devMode {
		// 加载.env文件
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}
	cl, err := GetClient()
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/models", HandleModels(cl))
	mux.HandleFunc("/v1/chat/completions", RelayHandler(cl, ChatCompletions))
	mux.HandleFunc("/v1/completions", RelayHandler(cl, Completions))
	mux.HandleFunc("/v1/embeddings", RelayHandler(cl, Embeddings))
	mux.HandleFunc("/v1/audio/speech", RelayHandler(cl, CreateSpeech))
	mux.HandleFunc("/v1/audio/transcriptions", RelayHandler(cl, Transcriptions))
	mux.HandleFunc("/v1/audio/translations", RelayHandler(cl, Translations))
	mux.HandleFunc("/v1/images/generations", RelayHandler(cl, CreateImage))
	mux.HandleFunc("/v1/images/edits", RelayHandler(cl, CreateImageEdit))
	mux.HandleFunc("/v1/images/variations", RelayHandler(cl, ImageVariations))

	handler := corsMiddleware(mux)

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", handler)
}

// CORS中间件
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头部
		w.Header().Set("Access-Control-Allow-Origin", getHeaderValue(r.Header, "Origin")) // 允许所有来源
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", getHeaderValue(r.Header, "Access-Control-Request-Method"))
			w.Header().Set("Access-Control-Allow-Headers", getHeaderValue(r.Header, "Access-Control-Request-Headers"))
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r) // 调用下一个处理程序
	})
}

func getHeaderValue(header http.Header, key string) string {
	values := header.Values(key)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}
