package tests

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestHttpProxyReq(t *testing.T) {
	proxyURL, err := url.Parse("sock5://127.0.0.1:7890") // 替换为你的代理地址和端口
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", "https://ipinfo.io", nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//if resp.StatusCode != http.StatusOK {
	//	panic("请求失败")
	//}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	println(string(body))
}
