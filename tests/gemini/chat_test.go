package gemini

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jiu-u/oai-adapter/api"
	"github.com/jiu-u/oai-adapter/clients/gemini"
	"io"
	"testing"
)

func TestChatParser(t *testing.T) {
	mediaMsg1 := []api.MediaContent{
		{
			Type: "image_url",
			ImageUrl: &api.MessageImageUrl{
				Url:    "https://example.com/image.jpg",
				Detail: "This is an example image.",
			},
		},
		{
			Type: "Text",
			Text: "hello",
		},
	}
	byteMsg1, err := json.Marshal(mediaMsg1)
	if err != nil {
		t.Fatal(err)
	}
	mediaMsg2 := []api.MediaContent{
		{
			Type: "input_audio",
			InputAudio: &api.MessageInputAudio{
				Data:   "b64hello",
				Format: "mp3",
			},
		},
		{
			Type: "Text",
			Text: "hello",
		},
	}
	byteMsg2, err := json.Marshal(mediaMsg2)
	if err != nil {
		t.Fatal(err)
	}
	imput := api.ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []api.Message{
			{Role: "developer", Content: []byte(`"hello"`)},
			{Role: "user", Content: []byte(`"hello"`)},
			{Role: "assistant", Content: byteMsg1},
			{Role: "user", Content: []byte(`"hello"`)},
			{Role: "user", Content: byteMsg2},
		},
	}
	client := gemini.Client{}
	resp, err := client.ConvertChatRequest(&imput)
	if err != nil {
		t.Fatal(err)
	}
	respStr, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(respStr))
}

func TestChatCompletions(t *testing.T) {
	chatReq := api.ChatRequest{
		Model: "gemini-1.5-flash",
		Messages: []api.Message{
			{Role: "user", Content: []byte(`"你好，这是一条测试消息"`)},
		},
		Stream: true,
	}
	resp, _, err := geminiClient.ChatCompletions(context.Background(), &chatReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Close()
	if !chatReq.Stream {
		data, err := io.ReadAll(resp)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(data))
	} else {
		scanner := bufio.NewScanner(resp)
		for scanner.Scan() {
			line := scanner.Text()
			// 过滤掉可能的心跳包
			if len(line) == 0 {
				continue
			}
			fmt.Println(line)
			fmt.Println()
		}

		// 检查流读取错误
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading stream: %s\n", err)
		}
	}
}
