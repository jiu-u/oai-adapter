package gemini

import (
	"context"
	"fmt"
	"github.com/jiu-u/oai-adapter/clients/gemini"
	"net/url"
	"os"
	"testing"
)

var geminiClient *gemini.Client

func TestMain(m *testing.M) {
	proxyURL, err := url.Parse("sock5://127.0.0.1:7890") // 替换为你的代理地址和端口
	if err != nil {
		panic(err)
	}
	geminiClient = gemini.NewClient("https://generativelanguage.googleapis.com",
		"AIxxxx-xxxxxxxxxxxxxxxxxxxxxxx",
		proxyURL)

	// 在所有测试运行之前执行的初始化代码
	fmt.Println("Before tests")

	// 运行测试
	code := m.Run()

	// 在所有测试运行之后执行的清理代码
	fmt.Println("After tests")

	// 退出程序
	os.Exit(code)
}

func TestGetModels(t *testing.T) {
	models, err := geminiClient.Models(context.Background())
	if err != nil {
		t.Errorf("Failed to get models: %v", err)
	}
	fmt.Println("Models:", models)
}
