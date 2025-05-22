package common

import (
	"net/http"
	"sync"
	"time"
)

var (
	defaultClient *http.Client
	initOnce      sync.Once
)

// initDefaultClient 初始化默认的 HTTP 客户端，包括从环境变量读取代理配置
func initDefaultClient() {
	transport := &http.Transport{
		//MaxIdleConns:        100,
		//MaxIdleConnsPerHost: 10,
		IdleConnTimeout: 90 * time.Second,
		Proxy:           http.ProxyFromEnvironment,
	}

	defaultClient = &http.Client{
		Transport: transport,
		Timeout:   15 * time.Minute,
	}
}

// SetDefaultClient 设置自定义的 HTTP 客户端
func SetDefaultClient(client *http.Client) {
	defaultClient = client
}

// GetDefaultClient 获取默认的 HTTP 客户端，如果还未初始化则先进行初始化
func GetDefaultClient() *http.Client {
	initOnce.Do(initDefaultClient)
	return defaultClient
}
